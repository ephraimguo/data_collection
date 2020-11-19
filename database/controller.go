package database

import (
	"data_collection/controllers"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func GetConn() (database *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root123_!@#"
	dbName := "data_collection"
	database, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	if err != nil {
		panic(err.Error())
	}

	return database // return a pointer of the database instance
}

func InsertSingleRecord(obj controllers.Insertable, db *sql.DB) {
	if obj.GetTimeStamp() == 0 {
		return
	}

	var err error
	var insert *sql.Stmt

	insert, err = db.Prepare("insert into deribit (timestamp, pair, bid_price, ask_price) value (?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	res, err := insert.Exec(obj.GetTimeStamp(), obj.GetPair(), obj.GetBidPrice(), obj.GetAskPrice())
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

	fmt.Printf("last row id: %v, num of rows affected: %v", rowId, rowAffected)
}
