package coinone

type JsonCurrency struct {
	Result string `json:"result"`
	ErrorCode string `json:"errorCode"`
	Currency string `json:"currency"`
	CurrencyType string `json:"currencyType"`
}

type JsonOrderbook struct {
	ErrorCode string `json:"errorCode"`
	Timestamp string `json:"timestamp"`
	Currency string `json:"currency"`
	Bid []struct {
		Price string `json:"price"`
		Qty string `json:"qty"`
	} `json:"bid"`
	Ask []struct {
		Price string `json:"price"`
		Qty string `json:"qty"`
	} `json:"ask"`
	Result string `json:"result"`
}
//Recent Complete Orders
type JsonTrades struct {
	ErrorCode string `json:"errorCode"`
	Timestamp string `json:"timestamp"`
	Result string `json:"result"`
	Currency string `json:"currency"`
	CompleteOrders []struct {
		Timestamp string `json:"timestamp"`
		Price string `json:"price"`
		Qty string `json:"qty"`
	} `json:"completeOrders"`
}

type JsonTicker struct {
	Volume float64 `json:"volume"`
	Last int64 `json:"last"`
	Timestamp int64 `json:"timestamp"`
	High int64 `json:"high"`
	Result string `json:"result"`
	ErrorCode string `json:"errorCode"`
	First float64 `json:"first"`
	Low int64 `json:"low"`
	Currency string `json:"currency"`
}