package deribit

import (
	"encoding/json"
)

func Parse(info []byte) *Resp {
	var resp Resp
	err := json.Unmarshal(info, &resp)
	if err != nil {
		panic(err)
	}

	if resp.Params.Quote.Timestamp == 0 {
		return nil
	}
	return &resp
}
