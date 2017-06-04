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

	templateBase, _ := template.ParseFiles("static_content/templates/index.html")
	templateBase.Execute(res, nil)
}
