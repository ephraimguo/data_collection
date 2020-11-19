package binance

import (
	"data_collection/controllers"
	"fmt"
)

type Resp struct {
	Quote Quote `json:"data"`
}

type Quote struct {
	controllers.Insertable
	Timestamp   int `json:"T"`
	Platform    string
	Pair        string     `json:"s"`
	BidPriceArr [][]string `json:"b"`
	AskPriceArr [][]string `json:"a"`
	BidPrice    float64
	AskPrice    float64
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

func (q Quote) ToString() string {
	return fmt.Sprintf("TimeStamp: %v, Platform: %v, Pair: %v, BidPrice: %v, AskPrice: %v",
		q.GetTimeStamp(), q.GetPlatform(), q.GetPair(), q.GetBidPrice(), q.GetAskPrice())
}
