package main

import (
	"time"
)

type Item struct {
	ID    string    `json:"id"`
	Added time.Time `json:"when-added"`
	Note  string    `json:"note"`
	Done  time.Time `json:"when-done"`
}
type Items []*Item

type Storage struct {
	items map[string]*Item
}

var storage Storage

func init() {
	storage.items = make(map[string]*Item)
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
