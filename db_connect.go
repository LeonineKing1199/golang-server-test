package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DBHandle is the MySQL database handle
var DBHandle *sql.DB

// InitDatabase establishes our MySQL connection
func InitDatabase() {
	user := "christian"
	password := "test123"

	db, err := sql.Open("mysql", user+":"+password+"@/dbname")

	if err == nil {
		DBHandle = db
	}

	log.Fatal(err)
}
