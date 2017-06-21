package poloniex

import (
	"strconv"
	"github.com/ironpark/coinex/trader"
	"time"
	"net/http"
	"github.com/ironpark/coinex/db"
	"encoding/json"
	"errors"
)

func NewTrader(key,secret string,pair string) *EX_Poloniex {
	c := &EX_Poloniex{
		apiKey:key,
		apiSecret:secret,
		httpClient:&http.Client{},
		crypto:pair,
		exchange:"poloniex",
	}
	dbClient,_ := db.Default()
	c.db = dbClient
	return c
}

func (c *EX_Poloniex) TickerData(resolution string) trader.TikerData{
	before := time.Now().Add(-time.Hour*24*30)
	data,_ := c.db.TradeHistory(c.crypto,"poloniex",before,time.Now().Add(time.Minute*60*3),DefaultLimit,resolution)
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

func (c *EX_Poloniex) Exchange() string{
	return c.exchange
}
