package poloniex

type poloOder struct {
	id int64
	price float64
	amount float64
	name string
}

func (oder poloOder) Cancel() error{
	return nil
}

func (oder poloOder) Price() float64{
	return oder.price
}

func (oder poloOder) Amount() float64{
	return oder.amount
}

func (oder poloOder) Name() string{
	return oder.name
}

func (oder poloOder) IsOpen() bool{
	return false
}
