package bucket

import (
	//log "github.com/sirupsen/logrus"
	cdb "github.com/ironpark/coinex/db"
	"github.com/asaskevich/EventBus"
	"sync"
	//"time"
	"github.com/ironpark/coinex/db"
	//"fmt"
	"github.com/ironpark/coinex/bucket/source"
)
const (
	/*support exchanges*/
	EXCHANGE_POLONIEX = "poloniex"
	EXCHANGE_BITTREX = "bittrex"
	EXCHANGE_UPBIT = "upbit"
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

type AssetStatus struct {
	First int64
	Last int64
	IsStop bool
}

type bucket struct {
	events EventBus.Bus
	db *cdb.CoinDB
	workers map[string]*Worker
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
			make(map[string]*Worker),
		}


		instance.workers[EXCHANGE_UPBIT] = NewWorker(bus, source.NewUpbit())
		//instance.workers = append(instance.workers,NewWorker(bus, source.NewPoloniex()))
		//instance.workers = append(instance.workers,NewWorker(bus, source.NewUpbit()))
		//instance.workers = append(instance.workers,NewWorker(bus, source.NewBitfinex()))

		bus.SubscribeAsync(TOPIC_DATA,instance.dataEvent,false)
	})
	return instance
}

func (bk *bucket)Status() map[string]map[cdb.Pair]string{
	return nil
}

func (bk *bucket)dataEvent(market string,pair cdb.Pair,data []cdb.MarketTake) {

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

func (bk *bucket)AddToTrack(market string,pair cdb.Pair) {
	bk.workers[market].Add(pair)
}

func (bk *bucket)Close()  {
	bk.events.Unsubscribe(TOPIC_DATA,instance.dataEvent)
}

func (bk *bucket)Run()  {
	//start All workers
	for _,e := range bk.workers{
		e.Run()
	}

}