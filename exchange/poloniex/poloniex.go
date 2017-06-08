package poloniex

import (
	"errors"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"strconv"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"net/url"
	//"fmt"
)

//JSON DataStruct
type EX_Poloniex struct {
	updateTarget []string
	updateCallback func (*EX_Poloniex,string)

	Prices map[string]float64
	apiKey string
	apiKeySecret string
	reqTime int64
	closeSignal chan bool

}

func NewEXPoloniex() (*EX_Poloniex) {
	c := &EX_Poloniex{
		reqTime: 1000*10,
		updateTarget:[]string{"BTC","ETH"},
		updateCallback:nil,
		Prices:map[string]float64{},
	}
	return c
}

func (ex *EX_Poloniex)SetApiKey(key,keySecret string) {
	ex.apiKey = key
	ex.apiKeySecret = keySecret
}
//public apis
func (ex *EX_Poloniex)GetTicker() ( map[string]JsonTickerData,error) {
	var jsonTicker map[string]JsonTickerData
	var err error = nil
	resp := ex.callPublicApi("returnTicker")
	if resp != nil {
		err = json.Unmarshal(resp, &jsonTicker)
		if err != nil{
			return jsonTicker,err
		}
	} else {
		errors.New("Fail Call Api")
	}
	return jsonTicker, nil
}

func (ex *EX_Poloniex)GetChartData(coin_name string,start,end time.Time,period uint64)(ChartData,error)  {
	var err error = nil

	command := "returnChartData"
	param :=command + "&currencyPair=" + coin_name +
		"&start=" + strconv.FormatUint(uint64(start.UTC().Unix()),10) +
		"&end=" + strconv.FormatUint(uint64(end.UTC().Unix()),10) +
		"&period=" + strconv.FormatUint(period,10)

	resp := ex.callPublicApi(param)
	chartData := ChartData{}
	if resp != nil {
		err = json.Unmarshal(resp, &chartData)
		if err != nil{
			return ChartData{},err
		}
	} else {
		errors.New("Fail Call Api")
	}
	return chartData, nil
}

func (ex *EX_Poloniex)GetTradeHistory(coin_name string,start,end time.Time)(TradeData,error){
//: https://poloniex.com/public?command=returnTradeHistory&currencyPair=BTC_NXT&start=1410158341&end=1410499372
	var err error = nil

	command := "returnTradeHistory"
	param :=command + "&currencyPair=" + coin_name +
		"&start=" + strconv.FormatUint(uint64(start.UTC().Unix()),10) +
		"&end=" + strconv.FormatUint(uint64(end.UTC().Unix()),10)

	resp := ex.callPublicApi(param)
	tradeData := TradeData{}
	if resp != nil {

		err = json.Unmarshal(resp, &tradeData)
		if err != nil{
			return TradeData{},err
		}
	} else {
		errors.New("Fail Call Api")
	}
	return tradeData, nil
}
//private(tradingApi) apis
//https://poloniex.com/tradingApi

func (ex *EX_Poloniex)GetMyBalances()(map[string]float64,error)  {
	var balances_json map[string]json.Number
	var balances map[string]float64 = make(map[string]float64)

	a := ex.callPrivateApi("returnBalances",nil)
	json.Unmarshal(a,&balances_json)

	for k, v := range balances_json {
		v_float ,_ := v.Float64()
		if v_float != 0{
			balances[k]=v_float
		}
	}
	return  balances,nil
}

func  (ex *EX_Poloniex) callPrivateApi(command string,parameters map[string]string) ([]byte) {
	client := http.Client{}

	form := url.Values{}
	form.Add("command", command)
	for key, value := range parameters {
		form.Add(key, value)
	}

	form.Add("nonce", strconv.FormatInt(time.Now().UnixNano(), 10))
	body := form.Encode()
	req, err := http.NewRequest("POST", "https://poloniex.com/tradingApi", strings.NewReader(body))
	if err != nil {
		return nil
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", ex.apiKey)
	mac := hmac.New(sha512.New, []byte(ex.apiKeySecret))
	mac.Write([]byte(body))
	signature := hex.EncodeToString(mac.Sum(nil))

	req.Header.Add("Sign", signature)
	resp, err := client.Do(req)

	if err != nil {
		return nil
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		//sometimes resp in "\\"
		contents = []byte(strings.Replace(string(contents), "\\", "", -1))
		return contents
	}
	return nil
}

func  (ex *EX_Poloniex) callPublicApi(command string) ([]byte) {
	client := http.Client{}
	client.Timeout = time.Second * 1000

	response, err := client.Get("https://poloniex.com/public?command=" + command)
	if err != nil {
		return nil
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err == nil {
		//sometimes resp in "\\"
		contents = []byte(strings.Replace(string(contents), "\\", "", -1))
		return contents
	}
	return nil
}

