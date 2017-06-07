package main

import "net/http"

// Send400 is a helper function to respond with a 400 status
func Send400(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.WriteHeader(http.StatusBadRequest)
	resWriter.Write([]byte("Invalid user request"))
	return
}

// Send500 is a helper function to respond with a 500 status
func Send500(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.WriteHeader(http.StatusInternalServerError)
	resWriter.Write([]byte("Internal server error"))
	return
}
