package main

import (
	"data_collection/database"
	"data_collection/utils"
)

func init() {

}

func main() {
	db := database.GetConn()

	go utils.Subscribe("deribit", db)
	go utils.Subscribe("binance", db)

	for {
	}
}
