package simulate

import (
	"fmt"
	"github.com/ironpark/coinex/exchange/poloniex"
	"time"
)

func LoadAllChartDataPoloniex(name string) (poloniex.ChartData,error){
	ex := poloniex.NewEXPoloniex()
	resp,err:=ex.GetChartData(name,time.Unix(1451606400,0),time.Now(),300)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return resp,nil
}

func LoadAllTradeDataPoloniex(name string) (poloniex.ChartData,error){
	ex := poloniex.NewEXPoloniex()
	resp,err:=ex.GetChartData(name,time.Unix(1451606400,0),time.Now(),300)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return resp,nil
}
