package bucket

import (
	"time"
	"log"
	"github.com/ironpark/coinex/db"
	"strings"
	"github.com/ironpark/go-poloniex"
)
type WorkerMessage struct {
	Type string
	Base string
	Pair string
	Exchange string
	Msg interface{}
}

type EventMessage struct {
	Type string
	Base string
	Pair string
	Exchange string
}

//EventListener Manage
type EventListener func(Ex,Pair,Type string,ListenerID int64)
//type EventListener interface {
//	Update(Type string)
//}
type AssetLisener struct {
	ID       int64
	Pairs    []string
	Listener EventListener
}
type ListenerStore map[string][]*AssetLisener

func (ls ListenerStore)Remove(ex string,id int64){

		for i,item:= range ls[ex] {
			if item.ID == id{
				ls[ex] = append(ls[ex][:i], ls[ex][i+1:]...)
			}
		}

}
func (ls ListenerStore)RemoveAll(){
	for k := range ls {
		ls[k] = []*AssetLisener{}
	}
}
func (ls ListenerStore)Add(ex string,pair []string,event EventListener)  {
	if ls[ex] == nil{
		ls[ex] = []*AssetLisener{}
	}
	if len(ls[ex]) == 0{
		ls[ex] = append(ls[ex],&AssetLisener{
			ID:0,
			Pairs:pair,
			Listener:event,
		})
	}else{
		ls[ex] = append(ls[ex],&AssetLisener{
			ID:int64(len(ls[ex])),
			Pairs:pair,
			Listener:event,
		})
	}
}
func (ls ListenerStore)Call(T,ex,pair string){
	if ls[ex] == nil{
		return
	}
	for _,item:= range ls[ex] {
		if contains(item.Pairs,pair){
			item.Listener(ex,pair,T,item.ID)
		}
	}
}
type bucket struct {
	Assets []*Target
	workers map[string]Worker
	msg chan WorkerMessage
	stop chan bool
	event chan EventMessage
	db *db.CoinDB
	listener ListenerStore
	globalListener []EventListener
}

type Target struct {
	Stop bool
	Exchange string
	Base string
	Pair string
	First time.Time
	Last time.Time
	Start time.Time
	End time.Time
	LastID int
}

func NewBucket() *bucket {
	db,_ := db.Default()
	bk := &bucket{
		Assets:[]*Target{},
		db:db,
		workers:make(map[string]Worker),
		msg:make(chan WorkerMessage),
		event:make(chan EventMessage),
		listener:make(ListenerStore),
		globalListener:[]EventListener{},
	}
	bk.workers["poloniex"] = Worker(NewPoloniexWorker(bk.msg,"",""))
	return bk
}

type TradeData interface {
	Type() string
	TradeID() int64
	Amount() float64
	Rate() float64
	Total() float64
	Date() time.Time
}
//if ex == "ALL" {
//bk.globalListener = listener
//return
//}
func (bk *bucket)AddGlobalEventListener(listener EventListener) int64{
	bk.globalListener = append(bk.globalListener,listener)
	return int64(len(bk.globalListener)-1)
}
func (bk *bucket)RemoveGlobalEventListener(index int64) {
	bk.globalListener = append(bk.globalListener[:index], bk.globalListener[index+1:]...)
}
func (bk *bucket)AddEventListener(ex string,listener EventListener,pairs... string) int64{
	bk.listener.Add(ex,pairs,listener)
	return int64(len(bk.listener[ex])-1)
}
func (bk *bucket)RemoveEventListener(ex string, id int64) {
	bk.listener.Remove(ex,id)
}

func (bk *bucket)Add(target *Target) {
	bk.Assets = append(bk.Assets,target)
	ex := target.Exchange
	ex = strings.ToLower(ex)
	bk.workers[ex].Add(target)

}

func (bk *bucket)Remove(target *Target) {
	ex := target.Exchange
	ex = strings.ToLower(ex)
	if bk.workers[ex] == nil{

	}else {
		bk.workers[ex].Remove(target)
	}
}

func (bk *bucket)insert(msg WorkerMessage)  {
	batch, _ := bk.db.NewBatchPoints()
	switch v :=  msg.Msg.(type) {
	case []poloniex.Trade:
		for _, trade := range v {
			point, _ := bk.db.NewPoint(
				db.Tags{
					"CP": msg.Base + "_" + msg.Pair,
					"base": msg.Base,
					"pair": msg.Pair,
					"ex":   msg.Exchange,
					"type": trade.Type,
				},
				db.Fields{
					"TradeID": trade.TradeID,
					"Amount":  trade.Amount,
					"Rate":    trade.Rate,
					"Total":   trade.Total,
				}, trade.Date.Time)
			batch.AddPoint(point)
		}
	}

	err := bk.db.Write(batch)
	if err != nil {
		log.Fatal(err)
	}
}

func (bk *bucket)work(){
	for {
		select {
		case msg:= <- bk.msg:
			if msg.Type == "Stop" {
				continue
			}
			bk.insert(msg)
			switch msg.Type {
			case "Load": //now stop
				bk.listener.Call(msg.Type,msg.Exchange,msg.Base+"_"+msg.Pair)
			case "Real": //real time
				bk.listener.Call(msg.Type,msg.Exchange,msg.Base+"_"+msg.Pair)
			}
			for i,item := range bk.globalListener {
				item(msg.Type,msg.Exchange,msg.Base+"_"+msg.Pair,int64(i))
			}

		case <- bk.stop:
			for _,worker := range bk.workers{
				worker.Stop() <- true
			}
			return
		}
	}
}

func (bk *bucket)Run()  {
	for _,worker := range bk.workers{
		go worker.Do()
	}
	go bk.work()
}

func (bk *bucket)Stop()  {
	bk.stop <- true
}