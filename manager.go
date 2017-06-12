package coinex

import (
	"github.com/ironpark/coinex/ex/poloniex"
	"github.com/ironpark/coinex/db"
	tr "github.com/ironpark/coinex/trader"
)


type Manager struct{
	traders []tr.Trader
	db *db.CoinDB
}

func NewManager()(*Manager){
	traders := []tr.Trader{}
	return &Manager{traders,nil}
}

func (ma *Manager) AddTrader(trader tr.Trader)  {
	ma.traders = append(ma.traders,trader)
}


func (ma *Manager) Start(){
	poloPairs := []string{}
	poloTraders := map[string][]tr.Trader{}
	for _,trader := range ma.traders{

			pair := trader.Pair()
			poloPairs = append(poloPairs,pair)
			poloTraders[pair] = append(poloTraders[pair] ,trader)

	}
	poloniex.PushApi(poloPairs, func(pair string,data []tr.TradeData) {
		for _,trader := range poloTraders[pair]{
			if len(data) == 0 {
				continue
			}
			trader.Call(trader,data[len(data)-1])
		}
	})
}