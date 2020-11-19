package controllers

type Insertable interface {
	GetTimeStamp() int
	GetPlatform() string
	GetPair() string
	GetBidPrice() float64
	GetAskPrice() float64
	ToString() string
}
