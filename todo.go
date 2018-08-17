package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ovlad32/hpg/todo"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

/*
func ItemsHandler(w http.ResponseWriter, r *http.Request) {
	var list TodoItems
	for _, v := range storage.items {
		list = append(list, v)
	}

	var e = json.NewEncoder(w)
	err := e.Encode(list)

	if err != nil {
		err = errors.Wrap(err, "")
		log.Println(err)
	}
	//w.Write([]byte("items"))
}
*/
func ItemAddHandler(w http.ResponseWriter, r *http.Request) {

	fx.Populate()
	//r.

}
func FirstBloodEP(
	logger *zap.Logger,
	h *mux.Router) error {
	h.HandleFunc("/fb",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("fitst blood"))
		},
	)
	return nil
}

func AllTodoItemsEP(
	logger *zap.Logger,
	todo *todo.Dispatcher,
	h *mux.Router) error {

	h.HandleFunc("/item/all",
		func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Got a request")

			items, err := todo.GetAll()
			data, err := json.Marshal(items)
			w.Header().Set("Content-Type", "application/json")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(data)
			}
		},
	)

	h.HandleFunc("/item/newid",
		func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Got a newid request")
			id, err := todo.NewId()
			type resp struct {
				Id string `json:"newid"`
			}
			vr := &resp{}
			vr.Id = id
			data, err := json.Marshal(vr)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
			w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
			w.Header().Set("Expires", "0")                                         // Proxies.
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(data)
			}
		},
	)

	return nil
}

//func (mux *http.ServeMux)
