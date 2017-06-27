package bucket

type bucket struct {
	workers []Worker
	running chan WorkerMessage
}

func NewBucket() *bucket {
	running := make(chan WorkerMessage)
	return &bucket{running:running}
}

func (bk *bucket)work(){
	for {
		select {
		case msg:= <- bk.running:
			switch msg.Type {
				//now stop
			case "Stop":
				//load backdatas
			case "Load":
				//realtime
			case "Real":
			}
			return
		default:

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

}