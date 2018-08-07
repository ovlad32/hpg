package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

//var app *fx.App

func main() {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	//, os.Kill

	//m.HandleFunc("/items", ItemsHandler).Methods("GET")
	//m.HandleFunc("/about", AboutHandler).Methods("GET")
	//m.HandleFunc("/item/add", ItemAddHandler).Methods("POST")
	/*
		hs := http.Server{
			Addr:     ":8090",
			ErrorLog: log.New(os.Stderr, "", 0),
			Handler:  m,
		}*/

	//http.ListenAndServe(":8090")
	/*go func() {
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
	*/
	app := fx.New(
		fx.Provide(provideCurrentPath),
		fx.Provide(NewZapLogger),
		fx.Provide(NewStorage),
		fx.Provide(NewHttpServer),
		fx.Invoke(FirstBloodEndpoint),
		fx.Invoke(AllTodoItemEndpoint),
	)
	_ = app
	app.Run()
}

func NewZapLogger() (result *zap.Logger, err error) {
	result, err = zap.NewDevelopment()
	return
}

type CurrentPath string

func provideCurrentPath() (result CurrentPath, err error) {
	s, err := os.Getwd()
	if err != nil {
		return
	}
	result = CurrentPath(s)
	return
}

func NewHttpServer(lc fx.Lifecycle, logger *zap.Logger) (result *mux.Router) {
	result = mux.NewRouter()
	hs := http.Server{
		Addr:     ":8090",
		ErrorLog: log.New(os.Stderr, "", 0),
		Handler:  result,
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) (err error) {
				log.Println("Server is running...")
				go func() {
					if err = hs.ListenAndServe(); err != nil {
						log.Fatal(err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) (err error) {
				hs.Shutdown(ctx)
				log.Println("Server stopped.")
				return nil
			},
		},
	)
	/*
		OnStart func(context.Context) error
	OnStop  func(context.Context) error*/
	return
}
