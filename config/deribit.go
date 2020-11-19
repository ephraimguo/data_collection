package config

var Deribit = map[string]string{
	"Scheme":  "wss",
	"Host":    "www.deribit.com",
	"Path":    "/ws/api/v2",
	"Payload": "{\"method\":\"public/subscribe\",\"params\":{\"channels\":[\"quote.BTC-PERPETUAL\",\"quote.ETH-PERPETUAL\"]},\"jsonrpc\":\"2.0\",\"id\":7}",
}
