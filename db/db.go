package db

import (
	"time"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/influxdata/influxdb/models"
	log "github.com/sirupsen/logrus"
	"fmt"
	"strings"
)

type Config struct{
	Host string
	UserName string
	Password string
}

type CoinDB struct {
	config Config
	client client.Client
	status int64
}

type Pair struct {
	Quote string
	Base string
}

type MarketTake struct {
	TradeID int
	Amount float64
	Rate float64
	Total float64
	Time time.Time
}

type OHLC struct {
	Open float64
	High float64
	Low float64
	Close float64
	Volume float64
	Time time.Time
}
//for influx db
type Tags map[string]string
type Fields map[string]interface{}

func Default()*CoinDB {
	return New(	Config{
		Host:     "",
		UserName: "",
		Password: "",
	})
}

func New(config Config) *CoinDB{
	coindb := &CoinDB{}
	var err error
	coindb.client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Host,
		Username: config.UserName,
		Password: config.Password,
	})
	if err != nil {
		log.Fatal("Influxdb connect error: ", err)
	}

	if _, _, err := coindb.client.Ping(30 * time.Second); err != nil {
		log.Fatal("Influxdb connect error: ", err)
	} else {
		coindb.status = 1
		log.Info("Influxdb connect successfully")
	}
	//TODO check database
	return coindb
}

//influx database wrapper
func (db *CoinDB)batch(dbname string)(client.BatchPoints, error){
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  dbname,
		Precision: "s",
	})
	return bp, err
}
//influx database wrapper
func (db *CoinDB)point(name string,tags Tags,fields Fields,date time.Time) (*client.Point,error){
	pt, err := client.NewPoint(
		name,
		tags,
		fields,
		date)
	return pt, err
}

//influx database wrapper
func (db *CoinDB) write(bp client.BatchPoints) (error) {
	if err := db.client.Write(bp); err != nil {
		return err
	}
	return nil
}
func (db *CoinDB) putMarket(market string) error {
	q := client.NewQuery(fmt.Sprintf(`CREATE DATABASE "market_%s"`, market), "", "")
	if response, err := db.client.Query(q); err != nil {
		log.Error(err)
		return err
	} else if err = response.Error(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
//Data to Point
func (db *CoinDB)ohlc(pair Pair,data OHLC) *client.Point{
	point, _ := db.point(pair.Base+"/"+pair.Quote,
		Tags{
			"pair": pair.Quote,
			"base": pair.Base,
		}, Fields{
			"Open":   data.Open,
			"High":   data.High,
			"Low":    data.Low,
			"Close":  data.Close,
			"Volume": data.Volume,
		}, data.Time)
	return point
}
func (db *CoinDB)take(pair Pair,data MarketTake) *client.Point{
	point, _ := db.point(pair.Base+"/"+pair.Quote,
		Tags{
			"Quote": pair.Quote,
			"Base": pair.Base,
		}, Fields{
			"TradeID": data.TradeID,
			"Amount":  data.Amount,
			"Rate":    data.Rate,
			"Total":   data.Total,
		}, data.Time)
	return point
}
//data insert
func (db *CoinDB)PutOHLC(market string,pair Pair,data OHLC,date time.Time) {
	batch, _ := db.batch(market)
	batch.AddPoint(db.ohlc(pair,data))
	db.write(batch)
}

func (db *CoinDB)PutOHLCs(market string,pair Pair,data... OHLC) {
	batch, _ := db.batch(market)
	for _,item := range data {
		batch.AddPoint(db.ohlc(pair,item))
	}
	db.write(batch)
}

func (db *CoinDB)PutMarketTake(market string,pair Pair,data MarketTake) {
	batch, _ := db.batch(market)
	batch.AddPoint(db.take(pair,data))
	db.write(batch)
}

func (db *CoinDB)PutMarketTakes(market string,pair Pair,data []MarketTake){
	batch, _ := db.batch(market)
	for _,item := range data {
		batch.AddPoint(db.take(pair,item))
	}
	db.write(batch)
}

// GetMarkets return the list of market name
func (driver *CoinDB) GetMarkets() (data []string) {
	data = []string{}
	q := client.NewQuery("SHOW DATABASES", "", "s")
	if response, err := driver.client.Query(q); err == nil && response.Err == "" && len(response.Results) > 0 {
		result := response.Results[0]
		if result.Err == "" && len(result.Series) > 0 && len(result.Series[0].Values) > 0 {
			for _, v := range result.Series[0].Values {
				if len(v) > 0 {
					name := fmt.Sprint(v[0])
					if strings.HasPrefix(name, "market_") {
						data = append(data, strings.TrimPrefix(name, "market_"))
					}
				}
			}
		}
	}
	return
}

// GetMarkets return the list of GetCurrencyPair from market
func (driver *CoinDB) GetCurrencyPairs(market string) (data []Pair) {
	data = []Pair{}
	q := client.NewQuery("SHOW MEASUREMENTS", "market_"+market, "s")
	if response, err := driver.client.Query(q); err == nil && response.Err == "" && len(response.Results) > 0 {
		result := response.Results[0]
		if result.Err == "" && len(result.Series) > 0 && len(result.Series[0].Values) > 0 {
			for _, v := range result.Series[0].Values {
				if len(v) > 0 {
					name := fmt.Sprint(v[0])
					spl := strings.Split(name,"/")
					data = append(data,Pair{
						Base:spl[0],
						Quote:spl[1],
					})
				}
			}
		}
	}
	return
}

// GetOHLC return the list of OHLC from CurrencyPair
func (driver *CoinDB) GetOHLC(market string, pair Pair,start,end time.Time,period time.Duration) (data []OHLC) {
	raw := fmt.Sprintf(
		`SELECT FIRST("Rate"), MAX("Rate"), MIN("Rate"), LAST("Rate"), SUM("Amount") FROM "%v" WHERE time >= %vs AND time < %vs GROUP BY time(%vs)`,
		pair.Base+"/"+pair.Quote, start, end, period)
	q := client.NewQuery(raw, "market_"+market, "s")
	data = []OHLC{}
	driver.execute(q, func(row models.Row) {
		for _,item := range row.Values {
			t, _ := time.Parse(time.RFC3339, item[0].(string))
			d := OHLC{
				Time:   t,
				Volume: item[5].(float64),
				Open:   item[1].(float64),
				High:   item[2].(float64),
				Low:    item[3].(float64),
				Close:  item[4].(float64),
			}
			data = append(data, d)
		}
	})
	return
}

func (driver *CoinDB) GetFirstDate(market string, pair Pair) time.Time {
	raw := fmt.Sprintf(`SELECT Amount FROM "%v" order by time asc limit 1`, pair.Base+"/"+pair.Quote)
	q := client.NewQuery(raw, "market_"+market, "s")
	var t time.Time = nil
	driver.execute(q, func(row models.Row) {
		t, _ = time.Parse(time.RFC3339, row.Values[0][0].(string))
	})
	return t
}

func (driver *CoinDB) GetLastDate(market string, pair Pair) time.Time {

	raw := fmt.Sprintf(`SELECT Amount FROM "%v" order by time desc limit 1`, pair.Base+"/"+pair.Quote)
	q := client.NewQuery(raw, "market_"+market, "s")
	t := time.Now().UTC()

	driver.execute(q, func(row models.Row) {
		t, _ = time.Parse(time.RFC3339, row.Values[0][0].(string))
	})

	return t
}

func (driver *CoinDB) execute(q client.Query,fn func(row models.Row)) {
	if response, err := driver.client.Query(q); err == nil && response.Err == "" && len(response.Results) > 0 {
		result := response.Results[0]
		if len(result.Series) > 0 {
			fn(result.Series[0])
		}
	}
}
