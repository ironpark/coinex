package source

import (
	"github.com/ironpark/coinex/db"
	"time"
)

type DataSource interface {
	Status() int //server status
	Interval() int
	BackData(pair db.Pair,date time.Time) ([]db.OHLC, bool)
	RealTime(pair db.Pair) ([]db.OHLC)
	Name() string //exchange name
	MarketCreated(pair db.Pair) time.Time
}

