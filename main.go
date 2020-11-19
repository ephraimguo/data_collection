package main

import (
	"data_collection/database"
	"data_collection/utils"
)

// var addr = flag.String("addr", "fstream.binance.com/stream?streams=btcusdt@depth10/ethusdt@depth10@500ms", "http service address")

func main() {
	db := database.GetConn()
	utils.Subscribe("deribit", db)
}
