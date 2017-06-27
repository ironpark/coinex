package bucket

import (
	"time"
	"github.com/ironpark/go-poloniex"
	"log"
	"github.com/ironpark/coinex/db"
)
type Pair struct {
	Base string
	Pair string
	First time.Time
	Last time.Time
	Start time.Time
	End time.Time
}
type Worker interface {
	Do()
}

type WorkerMessage struct {
	Type string
	Pair string
	Exchange string
	Msg interface{}
}

type PoloniexWorker struct {
	db *db.CoinDB
	pair []Pair
	client *poloniex.Poloniex
	msg chan WorkerMessage
	running bool
}

type BittrexWorker struct {
	Msg chan WorkerMessage
}


func (w *PoloniexWorker)Do(){
	for {
		if !w.running {
			w.msg <- WorkerMessage{Type:"Stop"}
			break
		}
		for _, currentPair := range w.pair {
			errCNT := 0
			//time.Now().Unix()
			for {
				pair := currentPair.Base + "_" + currentPair.Pair
				history, err := w.client.MarketHistory(pair, currentPair.Start, currentPair.End)
				if err != nil {
					errCNT++
					log.Println(err)
					time.Sleep(time.Second * 1)
					continue
				}

				batch, _ := w.db.NewBatchPoints()
				for _, trade := range history {
					point, _ := w.db.NewPoint(
						db.Tags{
							"cryptocurrency": pair,
							"ex":             "poloniex",
							"type":           trade.Type,
						},
						db.Fields{
							"TradeID": trade.TradeID,
							"Amount":  trade.Amount,
							"Rate":    trade.Rate,
							"Total":   trade.Total,
						}, time.Time(trade.Date))
					batch.AddPoint(point)
				}
				err = w.db.Write(batch)
				if err != nil {
					log.Fatal(err)
				}
				if len(history) == 0 {
					break
				}

				if len(history) < 50000 {
					w.msg <- WorkerMessage{Type:"Real"}
					break
				}else{
					w.msg <- WorkerMessage{Type:"Load"}
				}
				currentPair.Start = time.Time(history[len(history)].Date)
				time.Sleep(time.Millisecond * 500)
			}
			time.Sleep(time.Millisecond * 200)
		}
	}
}

