package bucket

import (
	"time"
	"github.com/ironpark/go-poloniex"
	"github.com/asaskevich/EventBus"
	cdb "github.com/ironpark/coinex/db"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/IronPark/coinex/db"
)

//Poloniex
type workerStatus struct{
	lastid int
	start time.Time
	end time.Time
	realtime bool
}
type PoloniexWorker struct {
	client *poloniex.Poloniex
	status map[string]workerStatus
	db *cdb.CoinDB
}

func Poloniex() *PoloniexWorker{
	return &PoloniexWorker{
		client:poloniex.New("",""),
		status:make(map[string]workerStatus),
		db:cdb.Default(),
	}
}
//first call once
func (w *PoloniexWorker)Init(){
	for _,item :=range w.db.GetCurrencyPairs(EXCHANGE_POLONIEX) {
		currencyPair := item.Quote + "_" + item.Base
		status := w.status[currencyPair]
		status.realtime = true
		status.lastid = -1
		status.start = w.db.GetLastDate(EXCHANGE_POLONIEX,item)
		status.end = time.Now().UTC().Add(time.Hour)
		w.status[currencyPair] = status
	}
}

func (w *PoloniexWorker)PairInit(pair cdb.Pair,t time.Time){
	currencyPair := pair.Quote + "_" + pair.Base
	status ,ok := w.status[currencyPair]
	if !ok {
		status.start = t
		status.end = time.Now().UTC().Add(time.Hour)
		status.realtime = false
		status.lastid = -1
	}
}

func (w *PoloniexWorker)Do(bus EventBus.Bus,pair cdb.Pair) {
	currencyPair := pair.Quote + "_" + pair.Base
	status ,ok := w.status[currencyPair]
	defer func() {
		//update status
		w.status[currencyPair] = status
	}()

	if !ok {
		status.realtime = true
		status.lastid = -1
		status.start = time.Now()
		status.end = time.Now().UTC().Add(time.Hour)
	}

	t1 := unixMilli()
	history, err := w.client.MarketHistory(currencyPair, status.start, status.end)
	t2 := unixMilli() - t1
	//prevent multiple calls in a short time.
	if t2 < 400 {
		time.Sleep(time.Millisecond * time.Duration(400 - t2))
	}

	if err != nil {
		log.Println(err)
		return
	}
	length := len(history)
	//no data ||
	if len(history) == 0 || status.lastid == history[0].GlobalTradeID {
		return
	}
	//last id update
	status.lastid = history[0].GlobalTradeID

	historys := make([]db.MarketTake,length)
	for i,item:= range history{
		historys[i].Time = item.Date.Time
		historys[i].Rate = item.Rate
		historys[i].TradeID = item.GlobalTradeID
		historys[i].Total = item.Total
		historys[i].Amount = item.Amount
	}

	//switch
	if length == 50000 {
		status.realtime = false
	}else{
		status.realtime = true
	}

	if status.realtime { // realtime
		status.start = historys[0].Time
		status.end = time.Now().UTC().Add(time.Hour)
	}else{  // load
		status.end = historys[length-1].Time
	}

	//send to bucket use eventBus
	bus.Publish(TOPIC_DATA,EXCHANGE_POLONIEX,pair,historys)
}

func (w *PoloniexWorker)Exchange() string{
	return EXCHANGE_POLONIEX
}

func (w *PoloniexWorker)GetStatus() string{
	b,e := json.Marshal(w.status)
	if e != nil{
		log.Errorf("%s worker GetStatus() %v",EXCHANGE_POLONIEX,e)
		return "{}"
	}
	return string(b)
}