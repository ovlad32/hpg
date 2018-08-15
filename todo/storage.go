package todo

import (
	"errors"
	"strings"
)

type storage struct {
	items map[string]*Item
}

func newStorage() (s *storage, err error) {
	s = new(storage)
	s.items = make(map[string]*Item)
	return
}

func (s *storage) put(i *Item) (err error) {
	if i == nil {
		err = errors.New("appended item is nil!")
	}

	if i.ID = strings.TrimSpace(i.ID); i.ID == "" {
		err = errors.New("appended item ID is empty!")
	}
	s.items[i.ID] = i
	return nil
}

func (s *storage) getAll() (result Items, err error) {
	result = make(Items, 0, len(s.items))
	for _, v := range s.items {
		result = append(result, v)
	}
	return
}
