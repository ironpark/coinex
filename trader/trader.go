package trader

import "time"
import talib "github.com/markcheno/go-talib"
//import "fmt"
 type TradeData struct {
	ID int64
	Type string
	Amount float64
	Price float64
	Total float64
	Date time.Time
}

type TikerData map[string]interface{}
type MyOders map[string]string
type MyBalances []Balance

type Balance interface {
	Currency() string
	Available() float64
	All() float64
}

type Oder interface {
	Cancel() error
	Price() float64
	Amount() float64
	Name() string
	IsOpen() bool
}

type CurrencyPair interface {
	SellOder(amount,price float64) Oder
	BuyOder(amount,price float64) Oder
	TickerData(resolution string) TikerData
}

type Trader interface {
	//info
	MyOpenOders() []Oder
	MyBalance() (MyBalances,error)
	Exchange() string
	//trade
	SetTradeCallback(func(trader Trader,pair CurrencyPair))
	Pair(pair string) CurrencyPair
	//callback
	Call(Trader,TradeData)
}

//balance
func (balance MyBalances) Size() int{
	return len(balance)
}

//ticker
type TickerValue []float64

func (ticker TikerData) Low() TickerValue{
	return ticker["low"].([]float64)
}

func (ticker TikerData) High() TickerValue{
	return ticker["high"].([]float64)
}

func (ticker TikerData) Open() TickerValue{
	return ticker["first"].([]float64)
}

func (ticker TikerData) Close() TickerValue{
	return ticker["last"].([]float64)
}

func (ticker TikerData) Volume() TickerValue{
	return ticker["volume"].([]float64)
}

func (ticker TikerData) Avg() TickerValue{
	return ticker["avg"].([]float64)
}

func (ticker TikerData) WeightedAvg() TickerValue{
	return ticker["avg-w"].([]float64)
}

func (ticker TikerData) Time() []int64{
	return ticker["date"].([]int64)
}

func (value TickerValue) Last () float64 {
	return value[len(value)-1]
}
func (value TickerValue) Before (offset int ) float64 {
	return value[len(value)-1 -offset]
}
//math
func (value TickerValue) Sma (period int) TickerValue {
	return talib.Sma(value,period)
}

//To-Do..
func (value TickerValue) Ema (period int) TickerValue {
	return talib.Ema(value,period)
}

func (value TickerValue) Size () int {
	return len(value)
}