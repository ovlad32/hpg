package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type TodoItem struct {
	ID    string     `json:"id"`
	Added *time.Time `json:"when-added,omitempty"`
	Note  string     `json:"note"`
	Done  *time.Time `json:"when-done,omitempty"`
}
type TodoItems []*TodoItem

type Storage struct {
	items map[string]*TodoItem
}

func NewStorage() (s *Storage, err error) {
	storage := new(Storage)
	storage.items = make(map[string]*TodoItem)

	var added = time.Now().Add(time.Duration(-30) * time.Second)
	var done = time.Now()
	storage.items["1"] = &TodoItem{
		ID:    "1",
		Note:  "Note1",
		Added: &added,
		Done:  &done,
	}
	added = time.Now().Add(time.Duration(-30) * time.Second)
	storage.items["2"] = &TodoItem{
		ID:    "2",
		Note:  "Note2",
		Added: &added,
	}
	return storage, nil
}
func (s *Storage) Append(i *TodoItem) (err error) {
	return nil
}

func (s *Storage) Len() int {
	return len(s.items)
}
func NewItemID() (id string, err error) {
	id = RandStringBytesMaskImpr(10)
	return
}

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
func FirstBloodEndpoint(
	logger *zap.Logger,
	storage *Storage,
	h *mux.Router) error {
	h.HandleFunc("/fb",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("fitst blood"))
		},
	)
	return nil
}

func AllTodoItemEndpoint(
	logger *zap.Logger,
	storage *Storage,
	h *mux.Router) error {

	h.HandleFunc("/item/all",
		func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Got a request")
			data, err := json.Marshal(storage.items)
			w.Header().Set("Content-Type", "application/json")
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
