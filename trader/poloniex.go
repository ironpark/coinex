package trader

import (
	"github.com/ironPark/go-poloniex"
)

type TPoloniex struct {
	client *poloniex.Poloniex
}

func PoloniexTrader(key,secret string) *TPoloniex {
	ptrader := &TPoloniex{}
	ptrader.client = poloniex.New(key, secret)
	return ptrader
}

func (c *TPoloniex) TickerData(resolution string) TikerData{
	//before := time.Now().Add(-time.Hour*24*30)
	//data,_ := c.db.TradeHistory(c.crypto,"poloniex",before,time.Now().Add(time.Minute*60*3),DefaultLimit,resolution)
	return nil
}

func (c *TPoloniex) MyOpenOders() []Oder{
	return nil
}

func (c *TPoloniex) MyBalance() (balance MyBalances,err error){
	balances,err :=  c.client.GetBalance()
	if err != nil{
		return
	}
	return []Balance([]poloBalance(balances)),nil
}

func (c *TPoloniex) SellOder(pair string,amount,price float64) Oder{
	return poloOder{}
}

func (c *TPoloniex) BuyOder(pair string,amount,price float64) Oder{
	return poloOder{}
}

func (c *TPoloniex) SetTradeCallback(callback func(trader Trader,data TradeData)){
	//c.callback = callback
}


func (c *TPoloniex) Call(trader Trader,data TradeData){
	//if c.callback != nil {
	//	c.callback(trader,data)
	//}
}

func (c *TPoloniex) Pair() string{
	//return c.crypto
}

func (c *TPoloniex) Exchange() string{
	//return c.exchange
}
