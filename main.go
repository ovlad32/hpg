package main

//gsw "github.com/linxGnu/goseaweedfs"
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/ovlad32/hpg/sparse"
	"github.com/ovlad32/hpg/todo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

//var app *fx.App
func sparseTest() {
	sb := sparse.New()
	sb.SetBit(0, true)
	sb.SetBit(3, true)
	sb.SetBit(31, true)
	sb.SetBit(332, true)
	sb.SetBit(337, true)
	sb.SetBit(437, true)
	sb.SetBit(438, true)
	fmt.Println(sb.Cardinality())
	for i := int32(0); i >= 0; i = sb.NextSetBit(i + 1) {
		fmt.Printf(">%v\n", i)
	}
}
func main() {
	sparseTest()
	return
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	/*	sw := gsw.NewSeaweed("http", "seaweed.master:9334", []string{"memory"}, 2*1024*1024, 5*time.Minute)

		cm, fp, flid, err := sw.UploadFile("main.go", "", "")
		fmt.Printf("ChunkManifest: %v\n", cm)
		fmt.Printf("FilePart: %v\n", fp)
		fmt.Printf("fileId: %v\n", flid)
		fmt.Printf("Error: %v\n", err) */
	os.Exit(1)

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
		fx.Provide(NewHttpServer),
		fx.Provide(NewTodo),
		fx.Invoke(FirstBloodEP),
		fx.Invoke(AllTodoItemsEP),
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
				logger.Info("Server is running...")
				go func() {
					if err = hs.ListenAndServe(); err != nil {
						logger.Panic("Failed to start http server", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) (err error) {
				hs.Shutdown(ctx)
				logger.Info("Server stopped.")
				return nil
			},
		},
	)
	/*
		OnStart func(context.Context) error
		OnStop  func(context.Context) error*/
	return
}

func NewTodo(logger *zap.Logger) (r *todo.Dispatcher, err error) {
	return todo.New(todo.Logger(logger))
}
