package deribit

type Resp struct {
	Params Params `json:"params"`
}

type Params struct {
	Quote Quote `json:"data"`
}

type Quote struct {
	Timestamp int `json:"timestamp"`
	Platform  string
	Pair      string  `json:"instrument_name"`
	BidPrice  float64 `json:"best_bid_price"`
	AskPrice  float64 `json:"best_ask_price"`
}

func (q Quote) GetTimeStamp() int {
	return q.Timestamp
}

func (q Quote) GetPlatform() string {
	return q.Platform
}

func (q Quote) GetPair() string {
	return q.Pair
}

func (q Quote) GetBidPrice() float64 {
	return q.BidPrice
}

func (q Quote) GetAskPrice() float64 {
	return q.AskPrice
}
