package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// UserData is a simple datatype that we'll use
// across the project
type UserData struct {
	firstName string
	lastName  string
}

// CreateUser adds a user to our database, using the global
// DBHandler as specified in db_connect.go
func CreateUser(addUser *sql.Stmt) func(http.ResponseWriter, *http.Request) {
	return func(resWriter http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)

		var userData UserData

		err := decoder.Decode(&userData)
		if err != nil {
			fmt.Println(err.Error())
			resWriter.WriteHeader(http.StatusBadRequest)
			resWriter.Write([]byte("Invalid data POST'd to server"))
			return
		}

		defer req.Body.Close()

		resWriter.WriteHeader(http.StatusOK)
	}
}
