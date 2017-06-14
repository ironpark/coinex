package coinex

import (

	"github.com/ironpark/coinex/db"
	tr "github.com/ironpark/coinex/trader"
	"log"
	"time"
	"github.com/ironpark/coinex/web"
	"github.com/ironpark/coinex/ex/poloniex"
	"net/http"
	"io/ioutil"
	"runtime"
	"path"
	"fmt"
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
		//fmt.Println(d.Date)
	}
	err := ma.db.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
}

func (ma *Manager) Start(port int64){
	//for poloniex
	go func() {
		poloPairs := []string{}
		p := poloniex.NewTrader("","","")
		now := time.Now()
		log.Println("poloniex trader lists ...")
		for pair := range ma.traders["poloniex"] {
			poloPairs = append(poloPairs, pair)
			his := p.TradeHistory(pair, now.Add(-time.Hour*24), now)
			if his != nil {
				ma.insertTradeData(pair, "poloniex", his)
			}
		}
		log.Println("poloniex.PushApi init ...")
		poloniex.PushApi(poloPairs, func(pair string, data []tr.TradeData) {
			//insert Data
			ma.insertTradeData(pair, "poloniex", data)
			for _, trader := range ma.traders["poloniex"][pair] {
				trader.Call(trader, data[len(data)-1])
			}

			before := time.Now().Add(-time.Hour*24)
			hd,_ := ma.db.TradeHistory("BTC_ETH","poloniex",before,time.Now(),2000,"5m")

			//log.Println(finaljson)
			ma.sse.Notifier <- []byte(ma.historyToJson(hd))
		})
	}()
	log.Println("start web server (server send event)")
	ma.webServerStart(port)
}

func (ma *Manager) historyToJson(hd tr.TikerData) string {
	dates   := hd.Time()
	closes  := hd.Last()
	opens   := hd.First()
	highs   := hd.High()
	lows    := hd.Low()
	volumes := hd.Volume()
	//fmt.Println(opens)
	finaljson := "["
	for i := range dates {
		finaljson += fmt.Sprintf(
			"{\"date\":%d,\"open\":%.9f,\"high\":%.9f,\"low\":%.9f,\"close\":%.9f,\"volume\":%.9f}",
			dates[i],opens[i],highs[i],lows[i],closes[i],volumes[i])
		if len(volumes) -1 !=  i {
			finaljson += ","
		}
	}
	finaljson += "]"
	return finaljson
}

func (ma *Manager) webServerStart(port int64){
	http.HandleFunc("/", func(w http.ResponseWriter,req *http.Request) {
		file, err := ioutil.ReadFile(packagePath()+"/web/index.html")
		w.Header().Set("Content-Type","text/html; charset=utf-8")
		w.Write(file)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	http.HandleFunc("/sse",ma.sse.ServeHTTP)
	http.ListenAndServe("localhost:"+fmt.Sprintf("%d",port),nil)
}

func packagePath() string{
	//package path
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	dir := path.Dir(filename)
	return dir
}