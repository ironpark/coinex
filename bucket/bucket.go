package bucket

type Bucket struct {
	workers []Worker
	running chan bool
}

func NewBucket() *Bucket {
	running := make(chan bool)
	return &Bucket{running:running}
}

func (bk *Bucket)work(){
	for {
		select {
		case <- bk.running:
			return
		default:
			bk.doWorkers()
		}
	}
}

func (bk *Bucket) doWorkers(){
	for _,worker := range bk.workers{

	}
}

func (bk *Bucket)Run()  {
	go bk.work()
}

func (bk *Bucket)Stop()  {

}