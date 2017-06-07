package main

import (
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// RootHandler is our most base handler and will simply render the homepage
func RootHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != "GET" {
		http.NotFound(res, req)
		return
	}

	tx, err := DBHandle.Begin()
	if err != nil {
		Send500(res, req)
		return
	}

	users := []UserData{}
	rows, err := tx.Query(`SELECT FirstName, LastName FROM UserData`)

	if err != nil {
		tx.Rollback()
		Send500(res, req)
		return
	}

	for rows.Next() {
		var (
			firstName, lastName string
		)

		// a scan error here implies corruption of the database data
		if err = rows.Scan(&firstName, &lastName); err != nil {
			tx.Rollback()
			Send500(res, req)
			return
		}

		users = append(users, UserData{FirstName: firstName, LastName: lastName})
	}

	rows.Close()

	tx.Commit()
	// decent protection against XSS
	res.Header().Set("Content-Security-Policy", "script-src 'self'")

	templateBase, _ := template.ParseFiles("static_content/templates/index.html")
	templateBase.Execute(res, struct{ Users []UserData }{Users: users})
}
