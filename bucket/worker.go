package bucket

import "time"
type Pair struct {
	Base string
	Pair string
	First time.Time
	Last time.Time
}

type WorkerMessage struct {
	Type string
	Message interface{}
}

type PoloniexWorker struct {
	pair []Pair
	Msg chan WorkerMessage
}

type BittrexWorker struct {
	Msg chan WorkerMessage
}

func (w *PoloniexWorker)Do(){
	//first

	//for _,item := range w.pair  {
	//	times := time.Now()
	//	start := item.First
	//
	//}
	//for  {
	//	item
	//
	//}
}

