package bithumb

//JSON DataStruct
//Public Apis

type JsonTicker struct {
	Data struct {
		AveragePrice string `json:"average_price"`
		ClosingPrice string `json:"closing_price"`
		Date         int    `json:"date"`
		MaxPrice     string `json:"max_price"`
		MinPrice     string `json:"min_price"`
		OpeningPrice string `json:"opening_price"`
		UnitsTraded  string `json:"units_traded"`
		Volume1day   string `json:"volume_1day"`
		Volume7day   string `json:"volume_7day"`
	} `json:"data"`
	Status string `json:"status"`
}

type JsonOrderbook struct {
	Data struct {
		Asks []struct {
			Price    string `json:"price"`
			Quantity string `json:"quantity"`
		} `json:"asks"`
		Bids []struct {
			Price    string `json:"price"`
			Quantity string `json:"quantity"`
		} `json:"bids"`
		OrderCurrency   string `json:"order_currency"`
		PaymentCurrency string `json:"payment_currency"`
		Timestamp       int    `json:"timestamp"`
	} `json:"data"`
	Status string `json:"status"`
}

type JsonRecentTransaction struct {
	Data []struct {
		Price           string `json:"price"`
		Total           string `json:"total"`
		TransactionDate string `json:"transaction_date"`
		Type            string `json:"type"`
		UnitsTraded     string `json:"units_traded"`
	} `json:"data"`
	Status string `json:"status"`
}

//Private Apis
type JsonAccount struct {
	Data struct {
		AccountID string `json:"account_id"`
		Balance   string `json:"balance"`
		Created   int    `json:"created"`
		TradeFee  string `json:"trade_fee"`
	} `json:"data"`
	Status string `json:"status"`
}


type JsonBalance struct {
	Data struct {
		AvailableBtc string `json:"available_btc"`
		AvailableKrw string `json:"available_krw"`
		InUseBtc     string `json:"in_use_btc"`
		InUseKrw     string `json:"in_use_krw"`
		MisuBtc      string `json:"misu_btc"`
		MisuDepoKrw  string `json:"misu_depo_krw"`
		MisuKrw      string `json:"misu_krw"`
		TotalBtc     string `json:"total_btc"`
		TotalKrw     string `json:"total_krw"`
		XcoinLast    string `json:"xcoin_last"`
	} `json:"data"`
	Status string `json:"status"`
}
