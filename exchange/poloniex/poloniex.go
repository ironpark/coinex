package poloniex

import (
	"strconv"
	"time"
	"log"
	"gopkg.in/jcelliott/turnpike.v2"
	"github.com/IronPark/coinex"
)
//JSON DataStruct

type EX_Poloniex struct {
	updateTarget []string
	updateCallback func (*EX_Poloniex,string)

	Prices map[string]float64
	apiKey string
	apiKeySecret string
	reqTime int64
	closeSignal chan bool

}

func NewEXPoloniex() (*EX_Poloniex) {
	c := &EX_Poloniex{
		reqTime: 1000*10,
		updateTarget:[]string{"BTC","ETH"},
		updateCallback:nil,
		Prices:map[string]float64{},
	}

	return c
}
func  (ex *EX_Poloniex) SetUpdateCall(callback func (*EX_Poloniex,string)) {
	ex.updateCallback = callback
}


func  (ex *EX_Poloniex) AutoUpdate() {
	var currencyPairOld float64
	go func() {
		for {
			c, err := turnpike.NewWebsocketClient(turnpike.JSON, "wss://api.poloniex.com/", nil, nil)
			c.ReceiveTimeout = 100 * time.Second;
			if err != nil {
				log.Fatal(err)
			}
			_, err = c.JoinRealm("realm1", nil)
			if err != nil {
				log.Fatal(err)
			}

			t := coinex.Microsectime()
			c.Subscribe("ticker", nil, func(args []interface{}, kwargs map[string]interface{}) {
				tag := args[0].(string)
				if tag == "BTC_XRP" {
					t = coinex.Microsectime()

					currencyPair, _ := strconv.ParseFloat(args[1].(string), 64)
					//last, _ := strconv.ParseFloat(args[2].(string), 64)
					//lowestAsk, _ := strconv.ParseFloat(args[2].(string), 64)
					//highestBid, _ := strconv.ParseFloat(args[3].(string), 64)
					if currencyPairOld == currencyPair {
						return
					}
				}
			})
			t2 :=  coinex.Microsectime()

			for {
				select {
				case <-ex.closeSignal:
					break
				default:
					if c == nil {
						break
					}
					//re-connection
					if coinex.Microsectime() - t2 >= 2000000 {
						t2 = coinex.Microsectime()
						c.Close()
						break
					}
				}
			}

		}
	}()
}


func  (ex *EX_Poloniex) callJsonApi(endpoint string, params string) ([]byte){

	return nil
}
