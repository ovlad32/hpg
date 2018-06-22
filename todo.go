package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Item struct {
	ID    string     `json:"id"`
	Added *time.Time `json:"when-added,omitempty"`
	Note  string     `json:"note"`
	Done  *time.Time `json:"when-done,omitempty"`
}
type Items []*Item

type Storage struct {
	items map[string]*Item
}

var storage Storage

func init() {
	storage.items = make(map[string]*Item)

	var added = time.Now().Add(time.Duration(-30) * time.Second)
	var done = time.Now()
	storage.items["1"] = &Item{
		ID:    "1",
		Note:  "Note1",
		Added: &added,
		Done:  &done,
	}
	added = time.Now().Add(time.Duration(-30) * time.Second)
	storage.items["2"] = &Item{
		ID:    "2",
		Note:  "Note2",
		Added: &added,
	}

}
func (s *Storage) Append(i *Item) (err error) {
	return nil
}

func (s *Storage) Len() int {
	return len(s.items)
}
func NewItemID() (id string, err error) {
	id = RandStringBytesMaskImpr(10)
	return
}

func ItemsHandler(w http.ResponseWriter, r *http.Request) {
	var list Items
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

func ItemAddHandler(w http.ResponseWriter, r *http.Request) {
	r.
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("fitst blood"))
}
