package trader

import "time"

type TradeData struct {
	ID int64
	Type string
	Amount float64
	Price float64
	Total float64
	Date time.Time
}

type TikerData map[string][]float64
type MyOders map[string]string
type MyBalances []Balance

type Balance struct {
	Name string
	Amount float64
	OnOders float64

}

type Oder interface {
	Cancel() error
	Price() float64
	Amount() float64
	Name() string
	IsOpen() bool
}

type Trader interface {
	TickerData(resolution string) TikerData
	//info
	MyOpenOders() []Oder
	MyBalance() (MyBalances,error)

	//trade
	SellOder(pair string,amount,price float64) Oder
	BuyOder(pair string,amount,price float64) Oder

	SetTradeCallback(func(trader Trader,data TradeData))
	Pair() string
	Call(Trader,TradeData)
	Exchange() string
}

//balance
func (balance MyBalances) Size() int{
	return len(balance)
}

//ticker
func (ticker TikerData) Low() []float64{
	return ticker["low"]
}

func (ticker TikerData) High()[]float64{
	return ticker["high"]
}

func (ticker TikerData) First()[]float64{
	return ticker["first"]
}

func (ticker TikerData) Last() []float64{
	return ticker["last"]
}

func (ticker TikerData) Volume()[]float64{
	return ticker["volume"]
}

func (ticker TikerData) Avg()[]float64{
	return ticker["avg"]
}

func (ticker TikerData) WeightedAvg()[]float64{
	return ticker["avg-w"]
}
