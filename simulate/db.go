package simulate

import (
	"fmt"
	"github.com/ironpark/coinex/exchange/poloniex"
	"github.com/influxdata/influxdb/client/v2"
	"time"
	"log"
)
//JSON DataStruct
type Simulater struct {
	addr string
	username string
	password string
	dbname string
	dbClient client.Client
}

func NewSimulater(addr,username,password,dbname string) (*Simulater,error) {
	s := &Simulater{
		addr: addr,
		username:username,
		password:password,
		dbname:dbname,
	}
	var err error = nil
	// Create a new HTTPClient
	s.dbClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
		Timeout: time.Minute*10,
	})
	if err != nil{
		return nil,err
	}
	return s,nil
}

func (db *Simulater)UpdateTradeHistory(name string) (error){
	if db.addr == ""{
		db.addr = "http://localhost:8086"
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db.dbname,
		Precision: "s",
	})

	if err != nil {
		return err
	}
	//Get Trading History
	pol := poloniex.NewEXPoloniex()
	start := time.Now()

	tradeHistorys, err := pol.GetTradeHistory(name, start.Add(-time.Hour*24*365), start)
	if err != nil {
		fmt.Println(err)
		return err
	}
	final := tradeHistorys[len(tradeHistorys)-1]


	count := 0
	layout := "2006-01-02 15:04:05"

	for {
		start, _ = time.Parse(layout, final.Date)

		trade_data, _ := pol.GetTradeHistory(name, start.Add(-time.Hour*24*365), start)
		for i := 0; i < len(trade_data); i++ {
			if final.TradeID == trade_data[i].TradeID {
				tradeHistorys = append(tradeHistorys, trade_data[i+1:]...)
				break
			}
			//fmt.Println(trade_data[i])
		}

		fmt.Println(final.Date, start,"COUNT",len(trade_data))
		if len(trade_data) < 50000{
			fmt.Println("load data break",count)
			break
		}
		count++
		final = trade_data[len(trade_data)-1]
	}
	for i := 0; i < len(tradeHistorys); i++ {
		item := tradeHistorys[i]

		// Create a point and add to batch
		tags := map[string]string{"cryptocurrency": name}
		fields := map[string]interface{}{
			"TradeID": item.TradeID,
			"GlobalTradeID":  item.GlobalTradeID,
			"Amount":  item.Amount,
			"Rate":    item.Rate,
			"Total":   item.Total,
			"Type":    item.Type,
		}
		date, _ := time.Parse(layout, item.Date)
		pt, err := client.NewPoint("TradeData", tags, fields, date)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}



	// Write the batch
	if err := db.dbClient.Write(bp); err != nil {
		log.Fatal(err)
	}
	return nil
}

func LoadAllChartDataPoloniex(name string) (poloniex.ChartData,error){
	ex := poloniex.NewEXPoloniex()
	resp,err:=ex.GetChartData(name,time.Unix(1451606400,0),time.Now(),300)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return resp,nil
}
