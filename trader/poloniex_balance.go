package trader

import "github.com/ironpark/go-poloniex"

type poloBalance poloniex.Balance

func (b poloBalance) Name()string  {
	return b.Currency
}
func (b poloBalance) All()float64  {
	return b.Balance
}

func (b poloBalance) OnOder()float64  {
	return b.Pending
}

func (b poloBalance) Free()float64  {
	return b.Available
}

func (b poloBalance) BTC()float64  {
	return b.Value
}

