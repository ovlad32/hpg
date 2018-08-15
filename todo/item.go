package todo

import "time"

type Item struct {
	ID    string     `json:"id"`
	Added *time.Time `json:"when-added,omitempty"`
	Note  string     `json:"note"`
	Done  *time.Time `json:"when-done,omitempty"`
}
type Items []*Item

func NewItemID() (id string, err error) {
	id = randStringBytesMaskImpr(10)
	return
}
