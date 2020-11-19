package database

import (
	"data_collection/config"
	"data_collection/controllers"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
)

func GetConn() *sql.DB {
	database, err := sql.Open(
		config.DBDriver,
		config.DBUser+":"+config.DBPass+"@/"+config.DBName)

	if err != nil {
		panic(err.Error())
	}

	return database // return a pointer of the database instance
}

func InsertSingleRecord(obj controllers.Insertable, db *sql.DB) {
	if obj.GetTimeStamp() == 0 {
		return
	}

	if exists(obj, db) {
		return
	}

	var err error
	var insert *sql.Stmt

	stmt := fmt.Sprintf("insert into %s (timestamp, pair, bid_price, ask_price) value (?, ?, ?, ?)", obj.GetPlatform())
	insert, err = db.Prepare(stmt)
	defer insert.Close()

	if err != nil {
		panic(err.Error())
	}

	res, err := insert.Exec(obj.GetTimeStamp(), obj.GetPair(), obj.GetBidPrice(), obj.GetAskPrice())
	if err != nil {
		err.Error()
		if e, ok := err.(*mysql.MySQLError); ok {
			if e.Number == 1062 {
				return
			}
		}
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
	fmt.Printf("platform: %v, %+v", obj.GetPlatform(), obj.ToString())
	fmt.Printf("\tlast row id: %v, num of rows affected: %v\n", rowId, rowAffected)
}

func exists(obj controllers.Insertable, db *sql.DB) bool {
	//stmt := fmt.Sprintf("select * from %v where timestamp = ? and pair = ?", obj.GetPlatform())
	//rows, err := db.Query(stmt, obj.GetTimeStamp(), obj.GetPair())
	//defer rows.Close()
	//if err != nil {
	//	if err != sql.ErrNoRows {
	//		log.Fatal(err.Error())
	//	}
	//
	//	return false
	//}
	//
	//return true
	var exist bool
	stmt := fmt.Sprintf("select * from %v where timestamp = ? and pair = ?", obj.GetPlatform())
	query := fmt.Sprintf("SELECT exists (%s)", stmt)
	err := db.QueryRow(query, obj.GetTimeStamp(), obj.GetPair()).Scan(&exist)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("error checking if row exists '%s' %v", stmt, err)
	}
	return exist
}
