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

type JsonTicker struct {
	BTCAMP JsonTickerData `json:"BTC_AMP"`
	BTCARDR JsonTickerData `json:"BTC_ARDR"`
	BTCBCN JsonTickerData `json:"BTC_BCN"`
	BTCBCY JsonTickerData `json:"BTC_BCY"`
	BTCBELA JsonTickerData `json:"BTC_BELA"`
	BTCBLK JsonTickerData `json:"BTC_BLK"`
	BTCBTCD JsonTickerData `json:"BTC_BTCD"`
	BTCBTM JsonTickerData`json:"BTC_BTM"`
	BTCBTS JsonTickerData`json:"BTC_BTS"`
	BTCBURST JsonTickerData `json:"BTC_BURST"`
	BTCCLAM JsonTickerData `json:"BTC_CLAM"`
	BTCDASH JsonTickerData `json:"BTC_DASH"`
	BTCDCR JsonTickerData `json:"BTC_DCR"`
	BTCDGB JsonTickerData `json:"BTC_DGB"`
	BTCDOGE JsonTickerData `json:"BTC_DOGE"`
	BTCEMC2 JsonTickerData `json:"BTC_EMC2"`
	BTCETC JsonTickerData `json:"BTC_ETC"`
	BTCETH JsonTickerData `json:"BTC_ETH"`
	BTCEXP JsonTickerData `json:"BTC_EXP"`
	BTCFCT JsonTickerData `json:"BTC_FCT"`
	BTCFLDC JsonTickerData `json:"BTC_FLDC"`
	BTCFLO JsonTickerData `json:"BTC_FLO"`
	BTCGAME JsonTickerData`json:"BTC_GAME"`
	BTCGNO JsonTickerData `json:"BTC_GNO"`
	BTCGNT JsonTickerData `json:"BTC_GNT"`
	BTCGRC JsonTickerData`json:"BTC_GRC"`
	BTCHUC JsonTickerData `json:"BTC_HUC"`
	BTCLBC JsonTickerData  `json:"BTC_LBC"`
	BTCLSK JsonTickerData  `json:"BTC_LSK"`
	BTCLTC JsonTickerData `json:"BTC_LTC"`
	BTCMAID JsonTickerData `json:"BTC_MAID"`
	BTCNAUT JsonTickerData `json:"BTC_NAUT"`
	BTCNAV JsonTickerData  `json:"BTC_NAV"`
	BTCNEOS JsonTickerData  `json:"BTC_NEOS"`
	BTCNMC JsonTickerData  `json:"BTC_NMC"`
	BTCNOTE JsonTickerData `json:"BTC_NOTE"`
	BTCNXC JsonTickerData `json:"BTC_NXC"`
	BTCNXT JsonTickerData `json:"BTC_NXT"`
	BTCOMNI JsonTickerData `json:"BTC_OMNI"`
	BTCPASC JsonTickerData `json:"BTC_PASC"`
	BTCPINK JsonTickerData `json:"BTC_PINK"`
	BTCPOT JsonTickerData `json:"BTC_POT"`
	BTCPPC JsonTickerData `json:"BTC_PPC"`
	BTCRADS JsonTickerData `json:"BTC_RADS"`
	BTCREP JsonTickerData `json:"BTC_REP"`
	BTCRIC JsonTickerData `json:"BTC_RIC"`
	BTCSBD JsonTickerData `json:"BTC_SBD"`
	BTCSC JsonTickerData   `json:"BTC_SC"`
	BTCSJCX JsonTickerData `json:"BTC_SJCX"`
	BTCSTEEM JsonTickerData `json:"BTC_STEEM"`
	BTCSTR JsonTickerData `json:"BTC_STR"`
	BTCSTRAT JsonTickerData `json:"BTC_STRAT"`
	BTCSYS JsonTickerData `json:"BTC_SYS"`
	BTCVIA JsonTickerData `json:"BTC_VIA"`
	BTCVRC JsonTickerData `json:"BTC_VRC"`
	BTCVTC JsonTickerData `json:"BTC_VTC"`
	BTCXBC JsonTickerData `json:"BTC_XBC"`
	BTCXCP JsonTickerData `json:"BTC_XCP"`
	BTCXEM JsonTickerData `json:"BTC_XEM"`
	BTCXMR JsonTickerData `json:"BTC_XMR"`
	BTCXPM JsonTickerData `json:"BTC_XPM"`
	BTCXRP JsonTickerData `json:"BTC_XRP"`
	BTCXVC JsonTickerData `json:"BTC_XVC"`
	BTCZEC JsonTickerData `json:"BTC_ZEC"`
	ETHETC JsonTickerData `json:"ETH_ETC"`
	ETHGNO JsonTickerData `json:"ETH_GNO"`
	ETHGNT JsonTickerData `json:"ETH_GNT"`
	ETHLSK JsonTickerData `json:"ETH_LSK"`
	ETHREP JsonTickerData `json:"ETH_REP"`
	ETHSTEEM JsonTickerData `json:"ETH_STEEM"`
	ETHZEC JsonTickerData `json:"ETH_ZEC"`
	USDTBTC JsonTickerData `json:"USDT_BTC"`
	USDTDASH JsonTickerData `json:"USDT_DASH"`
	USDTETC JsonTickerData `json:"USDT_ETC"`
	USDTETH JsonTickerData `json:"USDT_ETH"`
	USDTLTC JsonTickerData `json:"USDT_LTC"`
	USDTNXT JsonTickerData `json:"USDT_NXT"`
	USDTREP JsonTickerData `json:"USDT_REP"`
	USDTSTR JsonTickerData `json:"USDT_STR"`
	USDTXMR JsonTickerData `json:"USDT_XMR"`
	USDTXRP JsonTickerData `json:"USDT_XRP"`
	USDTZEC JsonTickerData `json:"USDT_ZEC"`
	XMRBCN JsonTickerData `json:"XMR_BCN"`
	XMRBLK JsonTickerData`json:"XMR_BLK"`
	XMRBTCD JsonTickerData `json:"XMR_BTCD"`
	XMRDASH JsonTickerData `json:"XMR_DASH"`
	XMRLTC JsonTickerData `json:"XMR_LTC"`
	XMRMAID JsonTickerData `json:"XMR_MAID"`
	XMRNXT JsonTickerData `json:"XMR_NXT"`
	XMRZEC JsonTickerData `json:"XMR_ZEC"`
}
