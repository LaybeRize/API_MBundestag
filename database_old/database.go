package database

import (
	"github.com/jmoiron/sqlx"
	"log"
)

var DB *sqlx.DB
var DataSource = "user=postgres password=root dbname=test sslmode=disable"
var Connected = false

func TestDatabase(tableDrop string, typeDrop string) {
	var err error
	if !Connected {
		DB, err = sqlx.Connect("postgres", DataSource)
		Connected = true
	}

	if err != nil {
		log.Fatal(err)
	}

	DB.MustExec(tableDrop)
	DB.MustExec(typeDrop)
}
