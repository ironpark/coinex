package coinex

import (

	"github.com/ironpark/coinex/db"
	tr "github.com/ironpark/coinex/trader"
	"log"
	"time"
	"github.com/ironpark/coinex/web"
	"net/http"
	"io/ioutil"
	"runtime"
	"path"
	"fmt"
	"github.com/ironpark/go-poloniex"
)


type Manager struct{
	sse *web.Broker
	traders map[string]map[string][]tr.Trader
	db *db.CoinDB
}

func NewManager()(*Manager){
	traders := map[string]map[string][]tr.Trader{}
	dbClient ,_:= db.Default()

	return &Manager{web.NewSSEServer(),traders,dbClient}
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
		//LastTradeHistory
		for pair := range ma.traders["poloniex"] {
			last,err := ma.db.LastTradeHistory(pair);
			if err != nil {
				fmt.Println("4day")
				last = last.Add(-time.Hour*24*4)
			}
			poloPairs = append(poloPairs, pair)
			his := p.TradeHistory(pair, last, now)
			if his != nil {
				ma.insertTradeData(pair, "poloniex", his)
			}
		}

		for  {
			time.Sleep(time.Second*1)
			for _,pair := range poloPairs{
				now := time.Now()
				bef := now.Add(-time.Minute)
				hds := p.TradeHistory(pair, bef, now)
				if len(hds) == 0{
					time.Sleep(time.Second*1)
					continue;
				}
				//fmt.Println(hds[len(hds)-1].Date.Local())
				ma.insertTradeData(pair, "poloniex", hds)
				//DB Data
				before := time.Now().Add(-time.Hour*24)
				before =  time.Date(before.Year(),before.Month(),before.Day(),0,0,0,0,time.UTC)

				for _,trader := range ma.traders["poloniex"][pair] {
					trader.Call(trader,hds[len(hds)-1])
				}

				hd,_ := ma.db.TradeHistory(pair,"poloniex",before,time.Now().Add(time.Hour*24),2000,"5m")
				ma.sse.Notifier <- []byte(ma.historyToJson(hd))

			}
		}

	}()
	log.Println("start web server (server send event)")
	ma.webServerStart(port)
}

func (ma *Manager) historyToJson(hd tr.TikerData) string {
	dates   := hd.Time()
	closes  := hd.Close()
	opens   := hd.Open()
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