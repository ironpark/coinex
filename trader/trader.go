package trader

import (
	talib "github.com/markcheno/go-talib"
	"github.com/ironpark/coinex/db"
)
//import "fmt"

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
	TickerData(resolution string)  db.TikerData
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
	Call(Trader, db.TradeData)
}

//balance
func (balance MyBalances) Size() int{
	return len(balance)
}
