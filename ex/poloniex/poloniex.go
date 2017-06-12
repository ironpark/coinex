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
	"crypto/tls"
	"net"
	"log"
)
const DefaultLimit = 500
//JSON DataStruct

type EX_Poloniex struct {
	apiKey string
	apiSecret string
	httpClient *http.Client
	callback func(trader.Trader,trader.TradeData)
	db *db.CoinDB
	crypto string
}

func NewTrader(key,secret string,pair string) *EX_Poloniex {
	c := &EX_Poloniex{
		apiKey:key,
		apiSecret:secret,
		httpClient:&http.Client{},
		crypto:pair,
	}
	return c
}

func (c *EX_Poloniex) TickerData(resolution string) trader.TikerData{
	before := time.Now().Add(-time.Hour*24*30)
	data,_ := c.db.TradeHistory("poloniex",c.crypto,before,time.Now(),DefaultLimit,resolution)
	return data
}

func (c *EX_Poloniex) MyOpenOders() []trader.Oder{
	return nil
}

func (c *EX_Poloniex) MyBalance() (balance trader.MyBalances,err error){
	balance = []trader.Balance{}
	r, err := c.do( "?command=returnCompleteBalances", "", false)
	if err != nil {
		return
	}
	response := make(map[string]interface{})
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}

	if response["error"] != nil {
		err = errors.New(response["error"].(string))
		return
	}

	for k, v := range response {
		values := v.(map[string]interface{})
		amount, _ := strconv.ParseFloat(values["available"].(string), 64)
		onOders, _ := strconv.ParseFloat(values["onOrders"].(string), 64)
		balance = append(balance,trader.Balance{
			Name:k,
			Amount:amount,
			OnOders:onOders,
		})
	}
	return
}

func (c *EX_Poloniex) SellOder(pair string,amount,price float64) trader.Oder{
	return poloOder{}
}

func (c *EX_Poloniex) BuyOder(pair string,amount,price float64) trader.Oder{
	return poloOder{}
}

func (c *EX_Poloniex) SetTradeCallback(callback func(trader trader.Trader,data trader.TradeData)){
	c.callback = callback
}


func (c *EX_Poloniex) Call(trader trader.Trader,data trader.TradeData){
	if c.callback != nil {
		c.callback(trader,data)
	}
}

func (c *EX_Poloniex) Pair() string{
	return c.crypto
}

//Base Code
func (ex *EX_Poloniex) do(ressource string, payload string, authNeeded bool) (response []byte, err error) {
	connectTimer := time.NewTimer(60 * time.Second)

	var rawurl string
	var method string

	if authNeeded {
		method = "POST"
		rawurl = "https://poloniex.com/tradingApi/"
	}else{
		method = "GET"
		rawurl = "http://poloniex.com/public/"
	}

	req, err := http.NewRequest(method, rawurl, strings.NewReader(payload))
	if err != nil {
		return
	}
	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

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

func getNewTrade(args []interface{}) []trader.TradeData{
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

func PushApi(pair []string,callback func(pair string,trades []trader.TradeData)){
	ws, err := turnpike.NewWebsocketClient(turnpike.JSON, "wss://api.poloniex.com", &tls.Config{}, net.Dial)
	ws.ReceiveTimeout = time.Second*60*100
	if err != nil {
		log.Fatal(err)
	}
	_,err = ws.JoinRealm("realm1",nil) //3
	if err != nil {
		log.Fatal(err)
	}

	ws.ReceiveDone = make(chan bool)
	for _,p := range pair {
		ws.Subscribe(p, nil, func(args []interface{}, kwargs map[string]interface{}) {
			if args == nil || len(args) == 0 {
				return
			}
			trades := getNewTrade(args)
			callback(p,trades)
		})
	}
	log.Println("listening for meta events")
	<-ws.ReceiveDone
}