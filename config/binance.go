package config

var Binance = map[string]string{
	"Scheme":     "wss",
	"Host":       "fstream.binance.com",
	"Path":       "stream",
	"ForceQuery": "true",
	"RawQuery":   "streams=btcusdt@depth10/ethusdt@depth10@500ms",
}
