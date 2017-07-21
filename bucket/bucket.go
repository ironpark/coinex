package bucket

import (
	log "github.com/sirupsen/logrus"
	cdb "github.com/ironpark/coinex/db"
	"github.com/asaskevich/EventBus"
	"sync"
	"time"
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

type Worker interface {
	Init()
	Do(EventBus.Bus, cdb.Pair)
	Exchange() string
	GetStatus() string
	PairInit(cdb.Pair,time.Time)
}

type Asset struct {
	Pair cdb.Pair
	First int64
	Last int64
}

type bucket struct {
	events EventBus.Bus
	db *cdb.CoinDB
	localAssets map[string][]Asset
	workers map[string]Worker
}
var instance *bucket

//singleton pattern
var once sync.Once

func GetInstance() *bucket {
	once.Do(func() {
		bus := EventBus.New()
		db:= cdb.Default()
		instance = &bucket{
			bus,
			db,
			make(map[string][]Asset),
			make(map[string]Worker),
		}
		bus.Subscribe("bucket:data",instance.dataEvent)
	})
	return instance
}
func (bk *bucket)dataEvent(market string,pair cdb.Pair,data []cdb.MarketTake) {
	//insert data in influxDB
	bk.db.PutMarketTakes(market,pair,data)

	bk.events.Publish(TOPIC_UPDATE,pair,data)
}

func (bk *bucket)UnSubscribe(topic string,event interface{}){
	bk.events.Unsubscribe(topic,event)
}

func (bk *bucket)Subscribe(topic string,event interface{}){
	bk.events.Subscribe(topic,event)
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

func (bk *bucket)AddToTrack(market string,pair cdb.Pair,t time.Time) {
	if bk.workers[market] == nil{
		log.Errorf("There is no worker process for the %s market",market)
	}

	bk.workers[market].PairInit(pair,t)
	as := Asset{pair,0,0, }
	//TODO set first ~ last date from database
	bk.localAssets[market] = append(bk.localAssets[market],as)
}


func (bk *bucket)Close()  {
	bk.events.Unsubscribe("bucket:data",instance.dataEvent)
}


func (bk *bucket)Run()  {
	wg := sync.WaitGroup{}
	for {
		for k, v := range bk.workers {
			market := k
			wg.Add(1)
			go func() {
				defer wg.Done()
				//assets
				for _, asset := range bk.localAssets[market] {
					v.Do(bk.events, asset.Pair)
				}
			}()
		}
		//Wait for All workers
		wg.Wait()
		// TODO if called Close -> break
	}
}