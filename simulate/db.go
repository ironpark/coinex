package simulate

import (
	"fmt"
	"github.com/ironpark/coinex/exchange/poloniex"
	"time"
)


func LoadAllDataPoloniex(name string) error{
	ex := poloniex.NewEXPoloniex()
	resp,err:=ex.GetChartData(name,1451606400,uint64(time.Now().UTC().Unix()),300)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("CNT",len(resp))
	return nil
}
