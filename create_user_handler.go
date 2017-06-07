package main

import (
	"database/sql"
	"encoding/json"
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
// Design rationale is that instead of simply re-rendering the entire HTML doc,
// we can instead just give the client a list of updated users so our server can
// function as an API that can be hit by clients that aren't browsers
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
			Send400(resWriter, req)
			return
		}

		defer req.Body.Close()

		// we make sure to escape the strings, very crucial
		// but our header's content security policy should be
		// our ultimate fail-safe
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
