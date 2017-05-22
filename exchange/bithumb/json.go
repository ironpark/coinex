package bithumb

//JSON DataStruct
type TickerJsonEXBitumb struct {
	Status string `json:"status"`
	Data TickerDataEXBitumb `json:"data"`
}

type TickerDataEXBitumb struct {
	OpeningPrice float64 `json:"opening_price,string"`
	ClosingPrice float64 `json:"closing_price,string"`
	Min_price float64 `json:"min_price,string"`
	Max_price float64 `json:"max_price,string"`
	Average_price  float64 `json:"average_price,string"`
	Units_traded float64 `json:"units_traded,string"`
	Volume1day float64 `json:"volume_1day,string"`
	Volume7day float64 `json:"volume_7day,string"`
	Date int64 `json:"date,int"`

}

// Account JSON structure
type AccountInfoEXBitumb  struct {
	Status string `json:"status"`
	Data AccountInfoDataEXBitumb `json:"data"`
}
type AccountInfoDataEXBitumb struct {
	Created int64 `json:"created,string"`
	Account_id string `json:"account_id"`
	Trade_fee float64 `json:"trade_fee,string"`
	Balance float64 `json:"balance,string"`
}
