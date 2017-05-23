package poloniex

import (
	"errors"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
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

func (ex *EX_Poloniex)GetTicker() (JsonTicker,error) {
	var err error = nil
	resp := ex.callPublicApi("returnTicker")
	ticker := JsonTicker{}
	if resp != nil {
		err = json.Unmarshal(resp, &ticker)
		if err != nil{
			return JsonTicker{},err
		}
	} else {
		errors.New("Fail Call Api")
	}
	return ticker, nil
}

func  (ex *EX_Poloniex) callPublicApi(command string) ([]byte) {
	response, err := http.Get("https://poloniex.com/public?command=" + command)
	if err != nil {
		return nil
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err == nil {
		//sometimes resp in "\\"
		contents = []byte(strings.Replace(string(contents), "\\", "", -1));
		return contents
	}
	return nil
}
