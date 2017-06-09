package simulate

import (
	"github.com/IronPark/coinex/model"
	"log"
	"time"
)

//JSON DataStruct
type tradeOder struct{
	Price float64
	Amount float64
	Status int32
}

type Trade struct {
	db *SimulaterDB
	sm *Simulater
}


type Simulater struct {
	oderID int64
	oders map[int64]tradeOder
	db     *SimulaterDB
	trade  Trade
}

func (tr *Trade)Sell(Price,Amount float64) int64 {
	
}

func (tr *Trade)Buy(Price,Amount float64)  int64{

	return tr.sm.newOder(Price,Amount)
}



func NewSimulater(addr,username,password,dbname string) (*Simulater,error) {
	database,_ := NewSimulaterDB(addr,username,password,dbname)

	s := &Simulater{
		db:database,
		oderID:0,
	}
	s.trade = Trade{
		db:database,
		sm:s,
	}

	return s,nil
}

func (sm *Simulater)newOder(Price,Amount float64) int64{
	sm.oderID++
	sm.oders[sm.oderID] = tradeOder{Price, Amount,0}
	//To-Do
	/*
	 if complite

	 */
	return sm.oderID
}

func (sm *Simulater)ModelTest(model model.TradeModel,coinPair string,start time.Time)  {
	name := model.Name()
	version := model.Version()
	log.Printf("Model Test Start [%s %s]",name,version)
}

func (sm *Simulater)AlphaModelTest(model model.TradeModel,coinPair string,start time.Time){

}
