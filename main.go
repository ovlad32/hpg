package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	//, os.Kill

	m := mux.NewRouter()
	m.HandleFunc("/items", ItemsHandler).Methods("GET")
	m.HandleFunc("/about", AboutHandler).Methods("GET")
	m.HandleFunc("/item/add", ItemAddHandler).Methods("POST")

	hs := http.Server{
		Addr:     ":8090",
		ErrorLog: log.New(os.Stderr, "", 0),
		Handler:  m,
	}

	//http.ListenAndServe(":8090")
	go func() {
		log.Println("Server is running...")
		if err := hs.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		<-signalChan
		log.Println("Got server shutdown signal...")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
		err := hs.Shutdown(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Server stopped.")
	}()

	runtime.Goexit()
}
