package bithumb

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	//"fmt"
)

type EX_Bitumb struct {
	updateTarget []string
	updateCallback func (*EX_Bitumb,string)

	Prices map[string]float64
	apiKey string
	apiKeySecret string
	reqTime int64
	closeSignal chan bool

}


func NewEXBitumb() (*EX_Bitumb) {
	c := &EX_Bitumb{
		reqTime: 10,
		updateTarget:[]string{"BTC","ETH"},
		updateCallback:nil,
		Prices:map[string]float64{},
	}

	return c
}
func  (ex *EX_Bitumb) SetUpdateCall(callback func (*EX_Bitumb,string)) {
	ex.updateCallback = callback
}

func  (ex *EX_Bitumb) tickerUpdate(order,payment string) bool {
	var ticker TickerJsonEXBitumb
	resp := ex.callJsonApi("/public/ticker/"+order, "")
	if resp != nil {
		json.Unmarshal(resp, &ticker)
	} else {
		return false
	}

	if int(ex.Prices[order]) == int(ticker.Data.ClosingPrice){
		return false
	}
	ex.Prices[order] = ticker.Data.ClosingPrice
	return true
}

func  (ex *EX_Bitumb) updatePrices() {

	for _,e := range ex.updateTarget {

		if ex.tickerUpdate(e,"KRW") {
			if ex.updateCallback != nil {
				ex.updateCallback(ex, e)
			}
		}

		// element is the element from someSlice for where we are
	}
}
func  (ex *EX_Bitumb) UpdateLoop(){
	for {
		time.Sleep(10 * time.Second)
		ex.updatePrices()
	}
}
func  (ex *EX_Bitumb) AutoUpdate() {
	ex.updatePrices()
	go ex.UpdateLoop()
}


func  (ex *EX_Bitumb) callJsonApi(endpoint string, params string) ([]byte){
	var api_url = "https://api.bithumb.com" + endpoint
	// Connects to Bithumb API server and returns JSON result value.
	response, err := http.Get(api_url)
	contents, err := ioutil.ReadAll(response.Body)
	if err == nil {
		return contents
	}
	return nil
}
