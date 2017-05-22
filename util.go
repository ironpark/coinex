package coinex

import (
	"time"
)

func Microsectime() int64 {
	return time.Now().UnixNano() /time.Millisecond.Nanoseconds()
}