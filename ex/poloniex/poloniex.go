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
	"fmt"
	"github.com/ironpark/coinex/trader"
	"github.com/ironpark/coinex/db"
	"gopkg.in/jcelliott/turnpike.v2"
	//"crypto/tls"
	"net"
	"log"
)
const DefaultLimit = 5000000
//JSON DataStruct

type EX_Poloniex struct {
	apiKey string
	apiSecret string
	httpClient *http.Client
	callback func(trader.Trader,trader.TradeData)
	db *db.CoinDB
	crypto string
	exchange string
}

func (ex *EX_Poloniex)TradeHistory(pair string,start time.Time,end time.Time) []trader.TradeData{
	var returnData []trader.TradeData = []trader.TradeData{}
	for {

		ressource := fmt.Sprintf("?command=returnTradeHistory&currencyPair=%s&start=%d&end=%d", pair, start.Unix(), end.Unix())
		result, err := ex.do(ressource, "", false)
		if err != nil {
			return nil
		}
		type trade map[string]interface{}
		trades := []trade{}
		err = json.Unmarshal(result, &trades)
		if err != nil {
			continue;
		}
		count := len(trades)
		if count == 0{
			continue;

		}

		data := make([]trader.TradeData, count)

		for i, x := range trades {
			TradeID := x["tradeID"].(float64)
			Amount, _ := strconv.ParseFloat(x["amount"].(string), 64)
			Rate, _ := strconv.ParseFloat(x["rate"].(string), 64)
			Total, _ := strconv.ParseFloat(x["total"].(string), 64)
			Type := x["type"].(string)
			Date, _ := time.Parse("2006-01-02 15:04:05", x["date"].(string))
			data[i].ID = int64(TradeID)
			data[i].Amount = Amount
			data[i].Price = Rate
			data[i].Total = Total
			data[i].Type = Type

			if err != nil {
				log.Panic(err)
			}
			data[i].Date = Date
			//fmt.Println(data[i],TradeID)
		}
		returnData = append(returnData,data...)

		first_d := data[0].Date.Format("2006-01-02 15:04:05")
		start_d := start.Format("2006-01-02 15:04:05")
		if len(data) == 50000 {

			time.Sleep(time.Second)
			fmt.Println("continue...",start.Format("2006-01-02 15:04:05")+"||",first_d+"||",start_d)
			end = data[len(data)-1].Date
			continue
		}
		break
	}

	return returnData
}

//Base Code
func (ex *EX_Poloniex) do(ressource string, payload string, authNeeded bool) (response []byte, err error) {
	connectTimer := time.NewTimer(60 * time.Second)

	var rawurl string
	var method string

	if authNeeded {
		method = "POST"
		rawurl = "https://poloniex.com/tradingApi" + ressource
	}else{
		method = "GET"
		rawurl = "http://poloniex.com/public" + ressource
	}

	req, err := http.NewRequest(method, rawurl, strings.NewReader(payload))
	if err != nil {
		return
	}
	if method == "POST" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	req.Header.Add("Accept", "application/json")

	if authNeeded {
		if len(ex.apiKey) == 0 || len(ex.apiSecret) == 0 {
			err = errors.New("You need to set API Key and API Secret to call this method")
			return
		}

		nonce := time.Now().UnixNano()
		req.Header.Add("Key", ex.apiKey)
		req.Form.Add("nonce", fmt.Sprintf("%d", nonce))
		mac := hmac.New(sha512.New, []byte(ex.apiSecret))
		_, err = mac.Write([]byte(req.URL.String()))
		sig := hex.EncodeToString(mac.Sum(nil))
		req.Header.Add("Sign", sig)
	}

	resp, err := ex.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
	}
	response = []byte(strings.Replace(string(response), "\\", "", -1))
	return response, err
}

func (ex *EX_Poloniex) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := ex.httpClient.Do(req)
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout on reading data from Bittrex API")
	}
}

func PushApi(pair []string,callback func(pair string,trades []trader.TradeData)){
	var ws *turnpike.Client
	var err error
	count := 0
	for ; ; {
		ws, err = turnpike.NewWebsocketClient(turnpike.JSON, "wss://api.poloniex.com",nil, net.Dial)
		if ws == nil {
			count++
			if count == 5 {
				log.Fatal(err)
			}
			log.Println(err,"retry ...",count)
			continue
		}
		break
	}

	ws.ReceiveTimeout = time.Second*60*100
	if err != nil {
		log.Fatal(err)
	}
	_,err = ws.JoinRealm("realm1",nil) //3
	if err != nil {
		log.Fatal(err)
	}
	getNewTrade := func (args []interface{}) []trader.TradeData {
		datas := []trader.TradeData{}
		for x := range args {
			data := args[x].(map[string]interface{})
			msgData := data["data"].(map[string]interface{})
			msgType := data["type"].(string)
			if msgType == "newTrade" {
				Type := msgData["type"].(string)
				TradeID, _ := strconv.ParseInt(msgData["tradeID"].(string), 10, 64)
				Rate, _ := strconv.ParseFloat(msgData["rate"].(string), 64)
				Amount, _ := strconv.ParseFloat(msgData["amount"].(string), 64)
				Total, _ := strconv.ParseFloat(msgData["total"].(string), 64)
				Date, _ := time.Parse("2006-01-02 15:04:05", msgData["date"].(string))
				datas = append(datas, trader.TradeData{
					Type:   Type,
					ID:     TradeID,
					Price:  Rate,
					Total:  Total,
					Amount: Amount,
					Date:   Date,
				})
			}
		}
		return datas
	}
	ws.ReceiveDone = make(chan bool)
	for x := range pair {
		currency := pair[x]
		log.Printf("poloniex push api '%s' subcribe\n",currency)
		ws.Subscribe(currency, nil, func(args []interface{}, kwargs map[string]interface{}) {
			if args == nil || len(args) == 0 {
				return
			}
			trades := getNewTrade(args)
			if len(trades) == 0 {
				return
			}
			callback(currency,trades)
			fmt.Println(trades[len(trades)-1].Date,trades[len(trades)-1].ID)
		})

	}
	log.Println("listening for meta events")
	<-ws.ReceiveDone
}