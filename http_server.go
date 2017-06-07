package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func makeMultiplexer() *http.ServeMux {

	serveMux := http.NewServeMux()

	rootHandler := http.HandlerFunc(RootHandler)
	fileHandler := http.FileServer(http.Dir("./static_content/"))

	addUser, err := DBHandle.Prepare(`
		INSERT INTO UserData (FirstName, LastName)
		VALUES (?, ?)
	`)

	defer addUser.Close()

	if err != nil {
		log.Fatal(err)
	}

	addUserHandler := http.HandlerFunc(CreateUser(addUser))

	serveMux.Handle("/users", addUserHandler)
	serveMux.Handle("/scripts/", fileHandler)
	serveMux.Handle("/", rootHandler)

	return serveMux
}

// StartServer is our main bootstrapping point, creating a living HTTP server
func StartServer() {

	// create a channel to receive OS signals
	// Interrupt is used for cross-platform reasons
	// we wish to fundamentally capture the Interrupts and then
	// gracefully release resources back to the system
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	InitDatabase()
	defer DBHandle.Close()

	serverAddress := ":8080"
	server := &http.Server{Addr: serverAddress, Handler: makeMultiplexer()}

	go func() {
		fmt.Println("starting server, listening at: " + serverAddress)
		fmt.Println("press ctrl-c to end")

		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// finally block until a signal is received
	<-stop

	// now we wish to "cancel the context" so to speak and release
	// all resources back to the OS
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)

	fmt.Println("\nServer closed and resources released!")
}
