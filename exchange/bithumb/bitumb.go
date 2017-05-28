package bithumb

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"errors"
)

type EX_Bitumb struct {
	apiKey string
}

func NewEXBitumb() (*EX_Bitumb) {
	c := &EX_Bitumb{
	}
	return c
}

func (ex *EX_Bitumb)GetTicker(currency string) (JsonTickerData,error) {
	var err error = nil
	resp := ex.callPublicApi("ticker/"+currency)
	ticker := JsonTicker{}
	if resp != nil {
		err = json.Unmarshal(resp, &ticker)
		if err != nil{
			return JsonTickerData{},err
		}
	} else {
		errors.New("Fail Call Api")
	}

	err = ex.getError(ticker.Status)
	if err != nil {
		return JsonTickerData{},err
	}

	return ticker.Data,nil
}

func (ex *EX_Bitumb) GetTransactions(currency string) ([]JsonRecentTransactionData,error) {
	var err error = nil
	transaction := JsonRecentTransaction{}
	resp := ex.callPublicApi("recent_transactions/"+currency)
	if resp != nil {
		err = json.Unmarshal(resp, &transaction)
		if err != nil{
			return nil,err
		}
	} else {
		errors.New("Fail Call Api")
	}

	err = ex.getError(transaction.Status)
	if err != nil {
		return nil,err
	}
	return transaction.Data,err
}

func (ex *EX_Bitumb)GetOrderbook(currency string) (orderbook JsonOrderbook,err error) {
	resp := ex.callPublicApi("orderbook/"+currency)
	if resp != nil {
		err = json.Unmarshal(resp, &orderbook)
		if err != nil{
			return JsonOrderbook{},err
		}
	} else {
		errors.New("Fail Call Api")
	}
	err = ex.getError(orderbook.Status)
	if err != nil {
		return JsonOrderbook{},err
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
	return err
}

func  (ex *EX_Bitumb) callPublicApi(endpoint string) ([]byte){
	var api_url = "https://api.bithumb.com/public/" + endpoint
	// Connects to Bithumb API server and returns JSON result value.
	response, err := http.Get(api_url)
	if err != nil {
		return nil
	}
	contents, err := ioutil.ReadAll(response.Body)

	if err == nil {
		//sometimes resp in "\\"
		contents = []byte(strings.Replace(string(contents),"\\","",-1));
		return contents
	}
	return nil
}
