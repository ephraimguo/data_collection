package binance

import (
	"encoding/json"
	"math"
	"strconv"
)

func Parse(info []byte) *Resp {
	var resp Resp
	err := json.Unmarshal(info, &resp)
	if err != nil {
		panic(err)
	}

	if resp.Quote.Timestamp == 0 {
		return nil
	}

	resp.Quote.BidPrice = getMax(resp.Quote.BidPriceArr)
	resp.Quote.AskPrice = getMin(resp.Quote.AskPriceArr)
	return &resp
}

func getMax(arr [][]string) float64 {
	var max float64
	for _, val := range arr {
		f, _ := strconv.ParseFloat(val[0], 64)
		if max < f {
			max = f
		}
	}

	return max
}

func getMin(arr [][]string) float64 {
	var min float64 = math.MaxInt64
	for _, val := range arr {
		f, _ := strconv.ParseFloat(val[0], 64)
		if min > f {
			min = f
		}
	}

	return min
}
