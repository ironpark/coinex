package trader

import "github.com/ironPark/go-poloniex"

type poloBalance poloniex.Balance

func (b poloBalance) Currency()string  {
	return b.Name
}
func (b poloBalance) Available()float64  {
	return b.Amount
}
func (b poloBalance) All()float64  {
	return b.OnOders + b.Amount
}

