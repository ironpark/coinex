package bucket

import (
	log "github.com/sirupsen/logrus"
	cdb "github.com/ironpark/coinex/db"
	"github.com/asaskevich/EventBus"
	"sync"
	"time"
	"github.com/ironpark/coinex/db"
	//"fmt"
)
const (
	/*support exchanges*/
	EXCHANGE_POLONIEX = "poloniex"
	EXCHANGE_BITTREX = "bittrex"
	EXCHANGE_BITFINEX = "bitfinex"

	/*subscribe topic list*/
	//This topic is called when data comes in
	TOPIC_DATA = "bucket:data"
	//This topic is called when db is updated.
	TOPIC_UPDATE = "bucket:update"
	TOPIC_TICKER = "bucket:ticker"
	TOPIC_STATUS = "bucket:status"

)
type UpdateEvent func(market string,pair cdb.Pair,data []cdb.MarketTake)
type TickerUpdateEvent func()
type StatusUpdateEvent func(market string,pair cdb.Pair,status AssetStatus)

type Worker interface {
	Init()
	Do(EventBus.Bus, cdb.Pair)
	Exchange() string
	GetStatus() string
	PairInit(cdb.Pair,time.Time)
}

type AssetStatus struct {
	First int64
	Last int64
	IsStop bool
}

type bucket struct {
	events EventBus.Bus
	db *cdb.CoinDB
	localAssets map[string]map[cdb.Pair]AssetStatus
	workers map[string]Worker
}
var instance *bucket

//singleton pattern
var once sync.Once

func Instance() *bucket {
	once.Do(func() {
		bus := EventBus.New()
		db:= cdb.Default()

		instance = &bucket{
			bus,
			db,
			make(map[string]map[cdb.Pair]AssetStatus),
			make(map[string]Worker),
		}

		//load from database
		for _,market :=range db.GetMarkets(){
			market_name := market
			if instance.localAssets[market] == nil {
				instance.localAssets[market] = make(map[cdb.Pair]AssetStatus)
			}

			mk := instance.localAssets[market]
			for _,pair :=range db.GetCurrencyPairs(EXCHANGE_POLONIEX) {
				status := mk[pair]
				status.Last = db.GetLastDate(market_name,pair).Unix()
				status.First = db.GetFirstDate(market,pair).Unix()

				status.IsStop = false //false
				instance.localAssets[market][pair] = status
			}
		}

		//poloniex worker (data crawler) add
		instance.AddWorker(Poloniex())
		bus.SubscribeAsync(TOPIC_DATA,instance.dataEvent,false)
	})
	return instance
}

func (bk *bucket)Status() map[string]map[cdb.Pair]AssetStatus{
	return bk.localAssets
}

func (bk *bucket)dataEvent(market string,pair cdb.Pair,data []cdb.MarketTake) {
	//insert data in influxDB
	bk.db.PutMarketTakes(market,pair,data)
	//log.Info(market,pair,data)
	state := bk.localAssets[market][pair]
	stateChange := false
	last := data[0].Time.Unix()
	if state.Last < last{
		state.Last = last
		stateChange = true
	}
	first := data[len(data)-1].Time.Unix()
	if state.First > first{
		state.First = first
		stateChange = true
	}

	if stateChange {
		bk.localAssets[market][pair] = state
		bk.events.Publish(TOPIC_STATUS, market, pair, state)
	}


	bk.events.Publish(TOPIC_UPDATE, market, pair, data)
	bk.events.Publish(TOPIC_TICKER + ":" + pair.ToString())

}

func (bk *bucket)unsubscribe(topic string,event interface{}){
	bk.events.Unsubscribe(topic,event)
}
func (bk *bucket)subscribe(topic string,event interface{}){
	bk.events.SubscribeAsync(topic,event,false)
}

//Subscribe
func (bk *bucket)SubTickerUpdate(pair db.Pair,fu TickerUpdateEvent){
	bk.subscribe(TOPIC_TICKER+":"+pair.ToString(),fu)
}
func (bk *bucket)SubStatusUpdate(fu StatusUpdateEvent){
	bk.subscribe(TOPIC_STATUS,fu)
}
func (bk *bucket)SubDBUpdate(fu UpdateEvent){
	bk.subscribe(TOPIC_UPDATE,fu)
}
//UnSubscribe
func (bk *bucket)UnSubTickerUpdate(pair db.Pair,fu TickerUpdateEvent){
	bk.unsubscribe(TOPIC_TICKER+":"+pair.ToString(),fu)
}
func (bk *bucket)UnSubStatusUpdate(fu StatusUpdateEvent){
	bk.unsubscribe(TOPIC_STATUS,fu)
}
func (bk *bucket)UnSubDBUpdate(fu UpdateEvent){
	bk.unsubscribe(TOPIC_UPDATE,fu)
}

func (bk *bucket)AddWorker(wk Worker) {
	if wk == nil {
		log.Panic("Worker is nil")
		return
	}
	market := wk.Exchange()
	if bk.workers[market] == nil {
		bk.workers[market] = wk
		bk.workers[market].Init()
	}else{
		log.Warnf("%s is an already added.",market)
	}
}
var lock = sync.RWMutex{}
func (bk *bucket)AddToTrack(market string,pair cdb.Pair,t time.Time) {
	lock.RLock()
	defer lock.RUnlock()
	if bk.workers[market] == nil{
		log.Errorf("There is no worker process for the %s market",market)
		return
	}
	log.Infof("add to track %s %s",market,pair.ToString())
	bk.workers[market].PairInit(pair,t)
	if bk.localAssets[market] == nil {
		bk.localAssets[market] = make(map[cdb.Pair]AssetStatus)
	}
	bk.localAssets[market][pair] = AssetStatus{First:99999999999,Last:-9999999999,IsStop:false}
	//TODO set first ~ last date from database
}

func (bk *bucket)StopToTrack(market string,pair cdb.Pair) {
	lock.RLock()
	defer lock.RUnlock()
	if bk.workers[market] == nil{
		log.Errorf("There is no worker process for the %s market",market)
		return
	}
	if bk.localAssets[market] == nil {
		return
	}
	status := bk.localAssets[market][pair]
	status.IsStop = true
	bk.localAssets[market][pair] = status
}

func (bk *bucket)Close()  {
	bk.events.Unsubscribe(TOPIC_DATA,instance.dataEvent)
}


func (bk *bucket)Run()  {
	wg := sync.WaitGroup{}
	for {
		for k, v := range bk.workers {
			market := k
			wg.Add(1)
			go func() {
				defer wg.Done()
				assetStatus:= bk.localAssets[market]
				//assets
				for pair := range bk.localAssets[market] {
					if assetStatus[pair].IsStop {
						continue
					}
					log.Info(pair.ToString())
					v.Do(bk.events, pair)
				}
			}()
		}
		//Wait for All workers
		wg.Wait()
		// TODO if called Close -> break
	}
}