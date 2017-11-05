package source

import (
	"github.com/ironpark/coinex/db"
	"fmt"
	"math/rand"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"sync"
	"strings"
)

type bittrexOHLC struct {
	Code string  `json:"code"`
	Utc  string  `json:"candleDateTime"`
	Kst  string  `json:"candleDateTimeKst"`
	O    float64 `json:"openingPrice"`
	H    float64 `json:"highPrice"`
	L    float64 `json:"lowPrice"`
	C    float64 `json:"tradePrice"`
	QV   float64 `json:"candleAccTradeVolume"`
	BV   float64 `json:"candleAccTradePrice"`
}

type Bittrex struct {
	httpClient  *http.Client
	once sync.Once
	marketCreated map[db.Pair]time.Time
}

func NewBittrex() *Bittrex {
	return &Bittrex{
		&http.Client{},
		sync.Once{},
		make(map[db.Pair]time.Time),
	}
}

func (wk *Bittrex)MarketCreated(pair db.Pair) time.Time {
	wk.once.Do(func() {
		resp,_ := (&http.Client{}).Get("https://bittrex.com/api/v1.1/public/getmarketsummaries")
		body,_ := ioutil.ReadAll(resp.Body)

		js := make(map[string]interface{})
		json.Unmarshal(body,&js)

		for _,e:=range js["result"].([]interface{}) {
			name := e.(map[string]interface{})["MarketName"].(string)
			created := e.(map[string]interface{})["Created"].(string)
			split := strings.Split(name,"-")
			base := split[0]
			quote := split[2]

			t,_ := time.Parse("2006-01-02T15:04:05",created)
			wk.marketCreated[db.Pair{quote,base}] = t
		}
	})

	return wk.marketCreated[pair]
}

func (wk *Bittrex)Status() int {
	return 200
}

//1 minute
func (wk *Bittrex)Interval() int {
	return 1
}

func (wk *Bittrex)BackData(pair db.Pair,date time.Time) ([]db.OHLC,bool) {
	data :=  wk.getData(pair, date.UTC(),1000)
	if len(data) < 1000{
		return data,true
	}
	return data,false
}

func (wk *Bittrex)RealTime(pair db.Pair) []db.OHLC {
	return wk.getData(pair, time.Now().UTC(), 2)
}

func (wk *Bittrex)Name() string {
	return "bittrex"
}

func (wk *Bittrex)getData(pair db.Pair,date time.Time,limit int) []db.OHLC {
	return wk.get(pair,1,date.Format("2006-01-02T15:04:05.000Z"),limit)
}

func (wk *Bittrex)get(pair db.Pair,minute int,date string,limit int) []db.OHLC {
	//"BTC-ETH"
	url := fmt.Sprintf("https://crix-api.upbit.com/v1/crix/candles/minutes/%d?code=CRIX.UPBIT.%s&count=%d&to=%s&ciqrandom=%d",
		minute, fmt.Sprintf("%s-%s",pair.Base,pair.Quote), limit, "2017-10-28T16:27:00.000Z", rand.Int())

	req ,err := http.NewRequest(http.MethodGet,url,nil)

	if err != nil{
		return nil
	}

	req.Header.Set("Host","crix-api.upbit.com")
	req.Header.Set("Origin","https://upbit.com")
	req.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")

	resp,err := wk.httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil{
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil
	}

	var ohcls []bittrexOHLC
	if json.Unmarshal(body,ohcls) != nil{
		return nil
	}

	ohclContainer := make([]db.OHLC, len(ohcls))

	for i:=0;i<len(ohcls) ;i++  {
		ohclContainer[i].Open = ohcls[i].O
		ohclContainer[i].High = ohcls[i].H
		ohclContainer[i].Close = ohcls[i].C
		ohclContainer[i].Low = ohcls[i].L
		ohclContainer[i].Volume = ohcls[i].BV
	}

	return ohclContainer
}

