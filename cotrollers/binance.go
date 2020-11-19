package controllers

import (
	"data_collection/models/binance"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Customer struct {
	CustomerId   int
	CustomerName string
	SSN          string
}

func GetConnection() (database *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "datacollection"
	database, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	if err != nil {
		panic(err.Error())
	}

	return database // return a pointer of the database instance
}

func InsertSingleRecord(obj binance.Quote) {
	db := GetConnection()

	var err error
	var insert *sql.Stmt

	insert, err = db.Prepare("insert into binance (timestamp , pair, bid_price, ask_price) value (?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	res, err := insert.Exec(obj.Timestamp, obj.Pair, obj.BidPrice, obj.AskPrice)
	if err != nil {
		panic(err.Error())
	}

	rowId, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	fmt.Printf("last row id: %v, num of rows affected: %v", rowId, rowAffected)
}

func Parse(info []byte) *binance.Resp{
	var resp binance.Resp
	err := json.Unmarshal(info, &resp)
	if err != nil {
		panic(err)
	}

	resp.Quote.BidPrice = getMax(resp.Quote.BidPriceArr)
	resp.Quote.AskPrice = getMin(resp.Quote.AskPriceArr)
	return &resp
}

func getMax(arr [][]string) float64{
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

//func ParseQuote (info []byte) *binance.Quote{
//	var quote *binance.Quote
//	json.Unmarshal(info, quote)
//
//	return quote
//}

//func ParseFundingRate (info []byte) *binance.FundingRate {
//	var funding_rate *binance.FundingRate
//	json.Unmarshal(info, funding_rate)
//
//	return funding_rate
//}

