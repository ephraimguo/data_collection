package binance

type Resp struct {
	Quote Quote `json:"data"`
}

type Quote struct {
	Timestamp int `json:"T"`
	Platform string
	Pair string `json:"s"`
	BidPriceArr [][]string `json:"b"`
	AskPriceArr [][]string `json:"a"`
	BidPrice float64
	AskPrice float64
}