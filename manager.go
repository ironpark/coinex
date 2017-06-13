package coinex

import (
	"github.com/ironpark/coinex/ex/poloniex"
	"github.com/ironpark/coinex/db"
	tr "github.com/ironpark/coinex/trader"
	"log"
	"time"
	"github.com/ironpark/coinex/web"
	"net/http"
	"io/ioutil"
	"runtime"
	"path"
)


type Manager struct{
	sse *web.Broker
	traders map[string]map[string][]tr.Trader
	db *db.CoinDB
}

func NewManager()(*Manager){
	traders := map[string]map[string][]tr.Trader{}
	dbClient ,_:= db.Default()

	return &Manager{web.NewSSEServer(),traders,dbClient,}
}

func (ma *Manager) AddTrader(trader tr.Trader)  {
	exname := trader.Exchange()
	pair := trader.Pair()
	log.Println(exname,pair)
	if ma.traders[exname] == nil{
		ma.traders[exname] = map[string][]tr.Trader{}
	}
	ma.traders[exname][pair] = append(ma.traders[exname][pair],trader)
}

func (ma *Manager) insertTradeData(pair string,ex string,data []tr.TradeData)  {
	bp,_ := ma.db.NewBatchPoints()
	for _,d := range data {
		point ,_:= ma.db.NewPoint(
			db.Tags{
				"cryptocurrency": pair,
				"ex":             ex,
				"type":           d.Type,
			},
			db.Fields{
				"TradeID": d.ID,
				"Amount":  d.Amount,
				"Rate":    d.Price,
				"Total":   d.Total,
			}, d.Date)
		bp.AddPoint(point)
	}
	ma.db.Write(bp)
}

func (ma *Manager) Start(){
	//poloniex
	poloPairs := []string{}
	p := poloniex.NewTrader("","","")
	now := time.Now()
	for pair := range ma.traders["poloniex"]{
		poloPairs = append(poloPairs,pair)
		his := p.TradeHistory(pair,now.Add(-time.Hour*24),now)
		if his != nil{
			ma.insertTradeData(pair,"poloniex",his)
		}
	}

	go func() {
		poloniex.PushApi(poloPairs, func(pair string, data []tr.TradeData) {
			//insert Data
			ma.insertTradeData(pair, "poloniex", data)
			for _, trader := range ma.traders["poloniex"][pair] {
				trader.Call(trader, data[len(data)-1])
			}
			ma.sse.Notifier <- []byte(pair)
		})
	}()
	//package path
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	dir := path.Dir(filename)


	http.HandleFunc("/trade", func(w http.ResponseWriter,req *http.Request) {
		file, err := ioutil.ReadFile(dir+"/web/index.html")
		w.Header().Set("Content-Type","text/html; charset=utf-8")
		w.Write(file)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	http.HandleFunc("/sse",ma.sse.ServeHTTP)

	http.ListenAndServe("localhost:3000",nil)
	//bittrex

}