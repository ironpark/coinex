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
// queryDB convenience function to query the database
func (db *Simulater) queryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: db.dbname,
	}
	if response, err := db.dbClient.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func (db *Simulater)insertTradeData(name string,insertData poloniex.TradeData) error{
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db.dbname,
		Precision: "s",
	})

	if err != nil {
		return err
	}
	layout := "2006-01-02 15:04:05"
	for i := 0; i < len(insertData); i++ {
		item := insertData[i]

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

func (db *Simulater)getTradeHistoryFirstDate(name string)time.Time{
	q := newQuery().From("TradeData").TAG("cryptocurrency",name).ASC("time").Limit(1).Build()
	res, err := db.queryDB(q)
	if err != nil {
		log.Fatal(err)
	}
	t, _ := time.Parse(time.RFC3339, res[0].Series[0].Values[0][0].(string))
	return t
}

func (db *Simulater)GetTradeHistory(name string,start,end time.Time,resolution time.Duration) (error){
	db.getTradeHistoryFirstDate(name)
	//get first time data
	//q := fmt.Sprintf("SELECT * FROM %s LIMIT %d", "TradeData",20)


	//q = fmt.Sprintf("SELECT * FROM %s LIMIT %d", "TradeData", 20)
	//res, err = db.queryDB( q)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for i, row := range res[0].Series[0].Values {
	//	t, err := time.Parse(time.RFC3339, row[0].(string))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	val := row[1].(string)
	//	log.Printf("[%2d] %s: %s\n", i, t.Format(time.Stamp), val)
	//}
	return nil
}

func (db *Simulater)UpdateTradeHistory(name string) (error){
	//Get Trading History
	pol := poloniex.NewEXPoloniex()
	start := time.Now()

	tradeHistorys, err := pol.GetTradeHistory(name, start.Add(-time.Hour*24*365), start)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//insert first data
	db.insertTradeData(name,tradeHistorys)

	final := tradeHistorys[len(tradeHistorys)-1]
	count := 0
	layout := "2006-01-02 15:04:05"

	for {
		start, _ = time.Parse(layout, final.Date)

		trade_data, err := pol.GetTradeHistory(name, start.Add(-time.Hour*24*365), start)
		if err != nil{
			break;
		}
		if len(trade_data) <= 0{
			break;
		}
		for i := 0; i < len(trade_data); i++ {
			if final.TradeID == trade_data[i].TradeID {
				db.insertTradeData(name,trade_data[i+1:])
				count+=len(trade_data)
				break
			}
		}


		fmt.Println("insert ",len(trade_data),"rows")
		final = trade_data[len(trade_data)-1]
		if len(trade_data) < 50000{
			fmt.Println("load data complete : ",count)
			break
		}
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
