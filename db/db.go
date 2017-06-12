package db

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"time"
	"log"
	"encoding/json"

)

type CoinDB struct {
	addr string
	username string
	password string
	dbname string
	dbClient client.Client
}

type DefaultConfig struct {
	Addr string
	Username string
	Password string
	DBname string
}

var (
	Config DefaultConfig
)

func init(){
	Config.Addr = "http://localhost:8086"
	Config.DBname = "coinex"
	Config.Username = ""
	Config.Password = ""
}

func SetUser(name,pass string)  {
	Config.Username = name
	Config.Password = pass
}

func Default() (*CoinDB,error) {
	return NewCoinDB(Config.Addr,Config.Username,Config.Password,Config.DBname)
}

func NewCoinDB(addr,username,password,dbname string) (*CoinDB,error) {

	s := &CoinDB{
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
func (db *CoinDB) queryDB(cmd string) (res []client.Result, err error) {
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


type Tags map[string]string
type Fields map[string]interface{}

func (db *CoinDB)NewBatchPoints()(client.BatchPoints, error){
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db.dbname,
		Precision: "s",
	})
	return bp, err
}

func (db *CoinDB) NewPoint(tags Tags,fields Fields,date time.Time) (*client.Point,error){
	pt, err := client.NewPoint(
		"TradeData",
		tags,
		fields,
		date)
	return pt, err
}

func (db *CoinDB) Write(bp client.BatchPoints) (error) {
	if err := db.dbClient.Write(bp); err != nil {
		return err
	}
	return nil
}


func (db *CoinDB)FirstTradeHistory(name string)time.Time{
	q := newQuery().From("TradeData").TAG("cryptocurrency",name).ASC("time").Limit(1).Build()
	res, err := db.queryDB(q)
	if err != nil {
		log.Fatal(err)
	}
	t, _ := time.Parse(time.RFC3339, res[0].Series[0].Values[0][0].(string))
	return t
}

func (db *CoinDB)LastTradeHistory(name string)time.Time{
	q := newQuery().From("TradeData").TAG("cryptocurrency",name).DESC("time").Limit(1).Build()
	res, err := db.queryDB(q)
	if err != nil {
		log.Fatal(err)
	}
	t, _ := time.Parse(time.RFC3339,res[0].Series[0].Values[0][0].(string))
	return t
}

func (db *CoinDB)getHistoryCount(name string) int64{
	q := newQuery().Select("count(TradeID)").From("TradeData").TAG("cryptocurrency",name).Build()
	res, err := db.queryDB(q)
	if err != nil {
		return 0
	}
	if len(res[0].Series) == 0{
		return 0
	}
	if len(res[0].Series[0].Values) == 0{
		return 0
	}
	if len(res[0].Series[0].Values[0]) == 0{
		return 0
	}

	count,_ := res[0].Series[0].Values[0][1].(json.Number).Int64()
	return count
}

func (db *CoinDB)TradeHistory(name, exchange string,start,end time.Time,limit int64,resolution string) (map[string][]float64,error) {

	q := newQuery().From("TradeData").TAG("cryptocurrency", name).TAG("ex", exchange).ASC("time").Limit(limit)
	q.Select(
		"MIN(Rate)",       // row
		"MAX(Rate)"   ,            // high
		"FIRST(Rate)" ,            // first(open)
		"LAST(Rate)"  ,            // last (close)
		"SUM(Total)"  ,            // volume
		"MEAN(Rate)"  ,            // Average
		"SUM(Total)/SUM(Amount)",  // weighted Average
		"STDDEV(Rate)",            // stddev
		"SPREAD(Rate)",            // diff between MIN MAX
	).GroupByTime(resolution).TIME(start, end).Build()
	fmt.Println(q.Build())
	res, err := db.queryDB(q.Build())
	if err != nil {
		log.Fatal(err)
	}


	result := res[0].Series[0].Values
	count := len(result)
	var output map[string][]float64 = make(map[string][]float64)
	output["low"]     = make([]float64, count)
	output["high"]    = make([]float64, count)
	output["first"]   = make([]float64, count)
	output["last"]    = make([]float64, count)
	output["volume"]  = make([]float64, count)
	output["avg"]     = make([]float64, count)
	output["avg-w"]   = make([]float64, count)
	output["stddev"]  = make([]float64, count)
	output["spread"]  = make([]float64, count)

	for i, row := range result {
		if row == nil {
			continue
		}

		t, err := time.Parse(time.RFC3339, row[0].(string))
		if err != nil {
			log.Fatal(err)
		}
		if row[1] == nil {
			continue
		}

		MIN, _ := row[1].(json.Number).Float64()
		MAX, _ := row[2].(json.Number).Float64()
		FIRST, _ := row[3].(json.Number).Float64()
		LAST, _ := row[4].(json.Number).Float64()
		VOLUME, _ := row[5].(json.Number).Float64()
		AVG, _ := row[6].(json.Number).Float64()
		AVGW, _ := row[7].(json.Number).Float64()

		var STDDEV,SPREAD float64
		if( row[8] != nil) {
			STDDEV, _ = row[8].(json.Number).Float64()
		}else{
			STDDEV = 0
		}

		if( row[9] != nil) {
			SPREAD, _ = row[9].(json.Number).Float64()
		}else{
			SPREAD = 0
		}
		output["spread"][i] = SPREAD
		output["stddev"][i] = STDDEV

		output["low"][i] = MIN
		output["high"][i] = MAX
		output["first"][i] = FIRST
		output["last"][i] = LAST
		output["volume"][i] = VOLUME
		output["avg"][i] = AVG
		output["avg-w"][i] = AVGW

		log.Printf("%s min %.8f max %.8f open %.8f close %.8f avg %.8f avg-w %.8f volume %.4f %f\n", t.Format(time.RFC3339), MIN, MAX, FIRST, LAST, AVG, AVGW,VOLUME,STDDEV)
	}

	return output,nil
}
//func (db *CoinDB)insetAllTradeHistory(name string,start time.Time) (error){
//	pol := ex.NewPoloniex()
//	tradeHistorys, err := pol.GetTradeHistory(name, start.Add(-time.Hour*24*365), start)
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	if(len(tradeHistorys) <= 0){
//		return nil
//	}
//	//insert first data
//	db.insertTradeData(name, tradeHistorys)
//
//	final := tradeHistorys[len(tradeHistorys)-1]
//	count := 0
//	layout := "2006-01-02 15:04:05"
//
//	for {
//		start, err = time.Parse(layout, final.Date)
//		if(err != nil){
//			fmt.Println(err)
//		}
//		trade_data, err := pol.GetTradeHistory(name, start.Add(-time.Hour*24*365), start)
//		if err != nil {
//			fmt.Println(err)
//			break;
//		}
//		if len(trade_data) <= 0 {
//			fmt.Println(len(trade_data))
//			break;
//		}
//		for i := 0; i < len(trade_data); i++ {
//			if final.TradeID == trade_data[i].TradeID {
//				db.insertTradeData(name, trade_data[i+1:])
//				count += len(trade_data)
//				break
//			}
//		}
//
//		final = trade_data[len(trade_data)-1]
//		fmt.Println("insert ", len(trade_data), "rows", final.Date)
//		if len(trade_data) < 50000 {
//			fmt.Println("load data complete : ", count)
//			break
//		}
//		time.Sleep(time.Second*1)
//	}
//	return nil
//}
//func (db *CoinDB)UpdateTradeHistory(name string) (error) {
//	//Get Trading History
//	pol := poloniex.NewEXPoloniex()
//	fmt.Println(db.getHistoryCount(name))
//	if db.getHistoryCount(name) == 0 {
//		fmt.Println("getHistoryCount")
//		db.insetAllTradeHistory(name, time.Now())
//	} else {
//		fmt.Println("Update")
//		//first Update
//		last := db.getTradeHistoryLastDate(name)
//		tradeHistorys, err := pol.GetTradeHistory(name, last, time.Now())
//		if err != nil {
//			fmt.Println(err)
//			return err
//		}
//		db.insertTradeData(name, tradeHistorys)
//		first := db.getTradeHistoryFirstDate(name)
//		db.insetAllTradeHistory(name, first)
//	}
//	return nil
//}
//
//func LoadAllChartDataPoloniex(name string) (poloniex.ChartData,error){
//	ex := poloniex.NewEXPoloniex()
//	resp,err:=ex.GetChartData(name,time.Unix(1451606400,0),time.Now(),300)
//	if err != nil {
//		fmt.Println(err)
//		return nil,err
//	}
//	return resp,nil
//}
