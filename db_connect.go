package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DBHandle is the MySQL database handle globally visible to
// the entire project
var DBHandle *sql.DB

func createDB(dbname, user, password string) {
	db, err := sql.Open("mysql", user+":"+password+"@/")
	defer db.Close()

	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE " + dbname)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE UserData (
			ID int NOT NULL AUTO_INCREMENT, 
			FirstName varchar(128), 
			LastName varchar(128),
			PRIMARY KEY (ID)
		)`)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully created database table: UserData")
}

// InitDatabase establishes our MySQL connection
func InitDatabase() {
	// customize these for your own SQL users and passwords
	user := "christian"
	password := "test123"
	dbname := "testinguserdatabase"

	db, err := sql.Open("mysql", user+":"+password+"@/"+dbname)

	if err != nil {
		fmt.Println("Error opening database connection")
		log.Fatal(err)
	}

	err = db.Ping()

	// with regards to error parsing, this is definitely a point of frailty and
	// should be subject to rigorous refactoring
	if err != nil && err.Error() == "Error 1049: Unknown database '"+dbname+"'" {
		db.Close()
		createDB(dbname, user, password)
		db, err = sql.Open("mysql", user+":"+password+"@/"+dbname)

		if err != nil {
			panic(err)
		}
	}

	db.Exec("USE " + dbname)
	db.Exec("SET SESSION TRANSACTION ISOLATION LEVEL REPEATABLE READ")

	DBHandle = db

	fmt.Println("Successfully established database connection and initialized global handle")
}
