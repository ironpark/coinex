package poloniex


type JsonTickerData struct {
	BaseVolume    string `json:"baseVolume"`
	High24hr      string `json:"high24hr"`
	HighestBid    string `json:"highestBid"`
	ID            int    `json:"id"`
	IsFrozen      string `json:"isFrozen"`
	Last          float64 `json:"last,string"`
	Low24hr       string `json:"low24hr"`
	LowestAsk     string `json:"lowestAsk"`
	PercentChange string `json:"percentChange"`
	QuoteVolume   string `json:"quoteVolume"`
}

type ChartData []struct {
	Date int64 `json:"date"`
	High float64 `json:"high"`
	Low float64 `json:"low"`
	Open float64 `json:"open"`
	Close float64 `json:"close"`
	Volume float64 `json:"volume"`
	QuoteVolume float64 `json:"quoteVolume"`
	WeightedAverage float64 `json:"weightedAverage"`
}

type TradeData []struct {
	GlobalTradeID int     `json:"globalTradeID"`
	TradeID       int     `json:"tradeID"`
	Date          string  `json:"date"`
	Type          string  `json:"type"`
	Rate          float64 `json:"rate,string"`
	Amount        float64 `json:"amount,string"`
	Total         float64 `json:"total,string"`
}