package models

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("mysql", "shiv:shiv123@tcp(localhost:3306)/taxify?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
}
