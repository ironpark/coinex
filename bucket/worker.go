package bucket

import (
	"time"
	"log"
	//"fmt"
	"github.com/ironpark/go-poloniex"
	"github.com/toorop/go-bittrex"
	"fmt"
)

type Worker interface {
	Do()
	Stop() chan bool
	Add(target *Target)
	Remove(target *Target)
	GetStatus()[]*Target
}
//Poloniex
type PoloniexWorker struct {
	pair []*Target
	client *poloniex.Poloniex
	msg chan WorkerMessage
	running chan bool
}

func NewPoloniexWorker(channel chan WorkerMessage,key, secret string) *PoloniexWorker {
	return &PoloniexWorker{
		pair:[]*Target{},
		client:poloniex.New(key,secret),
		msg:channel,
		running:make(chan bool),
	}
}

func (w *PoloniexWorker)Stop() chan bool{
	return w.running
}

func (w *PoloniexWorker)Add(target *Target){
	fmt.Println("Add1",target)
	for _,item := range w.pair {
		if item.Pair == target.Pair{
			if item.Base  == target.Base{
				return
			}
		}
	}
	fmt.Println("Add21",target)
	w.pair = append(w.pair,target)
}

func (w *PoloniexWorker)Remove(target *Target){
	for i,item := range w.pair {
		if item.Pair == target.Pair{
			if item.Base  == target.Base{
				w.pair = append(w.pair[:i], w.pair[i+1:]...)
				return
			}
		}
	}
}
func (w *PoloniexWorker)GetStatus()[]*Target{
	return w.pair
}
//Please note that making more than 6 calls per second to the public API (Poloniex)
func (w *PoloniexWorker)Do(){
	for {
		select {
		case <-w.running:
				w.msg <- WorkerMessage{Type: "Stop"}
			return
		default:
			for _, currentPair := range w.pair {
				if currentPair.Stop {
					continue
				}
				errCNT := 0
				pair := currentPair.Base + "_" + currentPair.Pair
				t1 := time.Now().Unix()
				history, err := w.client.MarketHistory(pair, currentPair.Start, currentPair.End)
				t2 := time.Now().Unix() - t1
				if t2 == 0 {
					time.Sleep(time.Millisecond * 1000)
				}

				if err != nil {
					errCNT++
					log.Println(err)
					time.Sleep(time.Second * 1)
					continue
				}
				if len(history) == 0 || currentPair.LastID == history[0].GlobalTradeID{
					continue
				}
				currentPair.LastID = history[0].GlobalTradeID

				fmt.Println(pair,history[len(history)-1].Date.Local(),"~",history[0].Date.Local(),len(history))
				currentPair.End = history[len(history)-1].Date.Time

				if currentPair.First.Unix() > currentPair.End.Unix() {
					currentPair.First = currentPair.End
				}

				if currentPair.Last.Unix() < history[0].Date.Time.Unix() {
					currentPair.Last = history[0].Date.Time
				}

				if len(history) < 50000 {
					w.msg <- WorkerMessage{Type: "Real", Base: currentPair.Base, Pair: currentPair.Pair, Msg: history,Exchange:"poloniex"}
					currentPair.Start = currentPair.Last
					currentPair.End = time.Now().UTC().Add(time.Hour)
					continue
				} else {
					w.msg <- WorkerMessage{Type: "Load", Base: currentPair.Base, Pair: currentPair.Pair, Msg: history,Exchange:"poloniex"}
				}

			}
		}
	}
}

type BittrexWorker struct {
	pair []*Target
	client *bittrex.Bittrex
	msg chan WorkerMessage
	running chan bool
}


func NewBittrexWorker(channel chan WorkerMessage,key, secret string) *BittrexWorker {
	return &BittrexWorker{
		pair:[]*Target{},
		client:bittrex.New(key,secret),
		msg:channel,
		running:make(chan bool),
	}
}

func (w *BittrexWorker)Stop() chan bool{
	return w.running
}

func (w *BittrexWorker)Add(target Target){
	for _,item := range w.pair {
		if item.Pair == target.Pair{
			if item.Base  == target.Base{
				return
			}
		}
	}
	w.pair = append(w.pair,&target)
}
//Please note that making more than 6 calls per second to the public API (Poloniex)
func (w *BittrexWorker)Do(){
	for {
		select {
		case <-w.running:
			w.msg <- WorkerMessage{Type: "Stop"}
			return
		default:

			for _, currentPair := range w.pair {
				if currentPair.Stop {
					continue
				}
			}
		}
	}
}
