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
	FirstName string
	LastName  string
}

func send400(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.WriteHeader(http.StatusBadRequest)
	resWriter.Write([]byte("Invalid data POST'd to server"))
	return
}

// CreateUser adds a user to our database, using the global
// DBHandler as specified in db_connect.go
func CreateUser(addUser *sql.Stmt) func(http.ResponseWriter, *http.Request) {
	return func(resWriter http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.NotFound(resWriter, req)
			return
		}

		decoder := json.NewDecoder(req.Body)

		var userData UserData

		err := decoder.Decode(&userData)
		if err != nil {
			fmt.Println(err.Error())
			send400(resWriter, req)
			return
		}

		defer req.Body.Close()

		firstName := userData.FirstName
		lastName := userData.LastName

		if len(firstName) == 0 || len(lastName) == 0 {
			send400(resWriter, req)
			return
		}

		resWriter.WriteHeader(http.StatusOK)
	}
}
