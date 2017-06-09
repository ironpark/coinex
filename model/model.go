package model

import (
	"github.com/IronPark/coinex/simulate"
	"time"
)

type TradeModel interface {
	//Model Info
	Name() string
	Version() string
	Run(sm *simulate.Trade,time time.Time)
	OderComplete(OderNum int64)
}

