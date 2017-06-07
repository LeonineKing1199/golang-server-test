package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

// UserData is a simple datatype that we'll use
// across the project
type UserData struct {
	FirstName string
	LastName  string
}

// CreateUser adds a user to our database, using the global
// DBHandler as specified in db_connect.go
// We add a user to the database and then re-fetch the whole thing and return
// it to the user so the front-end can re-render
// This is more friendly for single-page app design and not so much for
// multi-page apps which use formal form submissions followed by a page request
func CreateUser(addUser *sql.Stmt) func(http.ResponseWriter, *http.Request) {
	return func(resWriter http.ResponseWriter, req *http.Request) {

		// we wanna immediately reject all non-POST requests
		if req.Method != "POST" {
			http.NotFound(resWriter, req)
			return
		}

		// we then begin parsing the request body to see if we can
		// construct a valid piece of user data out of it
		decoder := json.NewDecoder(req.Body)

		var userData UserData

		err := decoder.Decode(&userData)
		if err != nil {
			fmt.Println(err.Error())
			Send400(resWriter, req)
			return
		}

		defer req.Body.Close()

		// we make sure to escape the strings, just to be extra safe
		firstName := html.EscapeString(userData.FirstName)
		lastName := html.EscapeString(userData.LastName)

		if len(firstName) == 0 || len(lastName) == 0 {
			Send400(resWriter, req)
			return
		}

		// we create a transaction that both writes to the database
		// and then reads the current contents of it to return to the user
		tx, err := DBHandle.Begin()
		if err != nil {
			Send500(resWriter, req)
			return
		}

		_, err = tx.Stmt(addUser).Exec(firstName, lastName)
		if err != nil {
			tx.Rollback()
			Send500(resWriter, req)
			return
		}

		rows, err := tx.Query(`SELECT FirstName, LastName FROM UserData`)
		if err != nil {
			tx.Rollback()
			Send500(resWriter, req)
			return
		}

		users := []UserData{}

		for rows.Next() {
			var (
				firstName, lastName string
			)

			// a scan error here implies corruption of the database data
			if err = rows.Scan(&firstName, &lastName); err != nil {
				tx.Rollback()
				Send500(resWriter, req)
				return
			}

			users = append(users, UserData{FirstName: firstName, LastName: lastName})
		}

		rows.Close()

		err = tx.Commit()
		if err != nil {
			fmt.Println(err.Error())
			tx.Rollback()
			Send500(resWriter, req)
			return
		}

		jsonResponse, err := json.Marshal(users)
		if err != nil {
			Send500(resWriter, req)
			return
		}

		resWriter.Header().Set("Content-Type", "application/json")
		resWriter.WriteHeader(http.StatusOK)
		resWriter.Write(jsonResponse)
	}
}
