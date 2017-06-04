package main

import (
	"log"
	"net/http"
)

func makeMultiplexer() *http.ServeMux {

	serveMux := http.NewServeMux()

	rootHandler := http.HandlerFunc(RootHandler)
	fileHandler := http.FileServer(http.Dir("./static_content/"))

	serveMux.Handle("/scripts/", fileHandler)
	serveMux.Handle("/", rootHandler)

	return serveMux
}

// StartServer is our main bootstrapping point, creating a living HTTP server
func StartServer() {
	InitDatabase()

	serverAddress := ":8080"
	log.Fatal(http.ListenAndServe(serverAddress, makeMultiplexer()))
}
