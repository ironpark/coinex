package bucket


import (
	"github.com/asaskevich/EventBus"
	coindb "github.com/ironpark/coinex/db"
	"github.com/google/go-cmp/cmp"
	"github.com/ironpark/coinex/bucket/source"
	"time"
	"os"
	"encoding/gob"
	"fmt"
)

type Status struct{
	First time.Time
	Last time.Time
	hasfirstTick bool
}

// Encode via Gob to file
func saveGob(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// Decode Gob file
func loadGob(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

type Worker struct {
	pairs []coindb.Pair
	// fisrt & last time data in local database
	status map[coindb.Pair]Status
	// is have the last historical data?
	inLast map[coindb.Pair]bool

	bus EventBus.Bus
	add chan coindb.Pair
	stop chan struct{}
	src source.DataSource
	db *coindb.CoinDB
}

func NewWorker(bus EventBus.Bus,src source.DataSource) *Worker{
	db := coindb.Default()

	worker := &Worker{
		pairs:db.GetCurrencyPairs(src.Name()),
		status:make(map[coindb.Pair]Status),
		bus:bus,
		add:make(chan coindb.Pair),
		stop:make(chan struct{}),
		src:src,
		db:db,
	}

	if worker.pairs == nil{
		worker.pairs = []coindb.Pair{}
	}

	for _,p := range worker.pairs {
		//get last/first update time
		worker.status[p] = Status{worker.db.GetLastDate(src.Name(), p),worker.db.GetFirstDate(worker.src.Name(), p),false}
	}

	var lastdata map[coindb.Pair]bool
	path := fmt.Sprintf("./%s",src.Name())
	if loadGob(path,lastdata) != nil {
		for _,p := range worker.pairs {
			lastdata[p] = false
		}
		saveGob(path,lastdata)
	}

	worker.inLast = lastdata
	return worker
}

func (w *Worker)Run(){
	path := fmt.Sprintf("./%s",w.src.Name())
	defer func(){
		close(w.add)
	}()
	for {
		select {
		default:
			for _,p := range w.pairs {
				//step : historical data
				status := w.status[p]
				if !w.inLast[p] {
					data,ok := w.src.BackData(p,status.First)
					status.First = data[len(data)-1].Time
					//update

					if ok {
						w.inLast[p] = true
						saveGob(path,w.inLast)
					}

					//insert database
					w.db.PutOHLCs(w.src.Name(),p,data...)
					w.status[p] = status
				}else{
					//step : synchronization data
					if time.Now().UTC().Unix() - status.Last.Unix() >= 120 {
						start := time.Now().UTC()

						for {
							data,_ := w.src.BackData(p,start)
							lastDataTime := data[len(data)-1].Time.Unix()
							status.Last = data[0].Time
							w.db.PutOHLCs(w.src.Name(),p,data...)
							w.status[p] = status

							if lastDataTime < w.status[p].Last.Unix(){
								break
							}
						}


					}else{
						//step : realtime data
						//TODO check the duplicated data
						data := w.src.RealTime(p)
						status.First = data[len(data)-1].Time

						w.db.PutOHLCs(w.src.Name(),p,data...)
						w.status[p] = status
						//TODO check point (really work?)
						w.bus.Publish(TOPIC_DATA, w.src.Name(), p, data)
					}
				}

			}
			//Add objects to track
		case pair := <-w.add:
			for p:= range w.pairs  {
				if cmp.Equal(p,pair) {
					continue
				}
			}
			w.pairs = append(w.pairs,pair)
		case <-w.stop:
			// stop
			return
		}
	}
}

func (w *Worker)Add(pair coindb.Pair){
	w.add <- pair
}

func (w *Worker)Stop(){
	close(w.stop)
}