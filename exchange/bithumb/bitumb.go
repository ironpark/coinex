package bithumb

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	//"fmt"
	"github.com/iris-contrib/errors"
)

type EX_Bitumb struct {
	Prices map[string]float64
	apiKey string
	apiKeySecret string
	reqTime int64
	closeSignal chan bool

}

func NewEXBitumb() (*EX_Bitumb) {
	c := &EX_Bitumb{
		reqTime: 10,
		Prices:map[string]float64{},
	}

	return c
}

func (ex *EX_Bitumb)GetTicker(currency string) (ticker JsonTicker,err error) {
	resp := ex.callPublicApi("ticker/"+currency)
	if resp != nil {
		json.Unmarshal(resp, &ticker)
	} else {
		errors.New("Fail Call Api")
	}
	err = ex.getError(ticker.Status)
	if err != nil {
		return nil,err
	}
	return ticker,err
}

func (ex *EX_Bitumb)GetTransactions(currency string) (transaction JsonRecentTransaction,err error) {
	resp := ex.callPublicApi("recent_transactions/"+currency)
	if resp != nil {
		json.Unmarshal(resp, &transaction)
	} else {
		errors.New("Fail Call Api")
	}
	err = ex.getError(transaction.Status)
	if err != nil {
		return nil,err
	}
	return transaction,err
}

func (ex *EX_Bitumb)GetOrderbook(currency string) (orderbook JsonOrderbook,err error) {
	resp := ex.callPublicApi("orderbook/"+currency)
	if resp != nil {
		json.Unmarshal(resp, &orderbook)
	} else {
		errors.New("Fail Call Api")
	}
	err = ex.getError(orderbook.Status)
	if err != nil {
		return nil,err
	}
	return orderbook,err
}


func (ex*EX_Bitumb) getError(status string)(err error) {
	switch status {
	case "5100":
		err = errors.New("Bad Request")
	case "5200":
		err = errors.New("Not Member")
	case "5300":
		err = errors.New("Invalid Apikey")
	case "5302":
		err = errors.New("Method Not Allowed")
	case "5400":
		err = errors.New("Database Fail")
	case "5500":
		err = errors.New("Invalid Parameter")
	case "5600":
		err = errors.New("CUSTOM NOTICE")
	case "5900":
		err = errors.New("Unknown Error")
	default:
		return nil
	}
}

func  (ex *EX_Bitumb) callPublicApi(endpoint string) ([]byte){
	var api_url = "https://api.bithumb.com/public/" + endpoint
	// Connects to Bithumb API server and returns JSON result value.
	response, err := http.Get(api_url)
	contents, err := ioutil.ReadAll(response.Body)
	if err == nil {
		return contents
	}
	return nil
}
