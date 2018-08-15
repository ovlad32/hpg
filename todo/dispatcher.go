package todo

import (
	"go.uber.org/zap"
)

type Option func(self *Dispatcher)

type Dispatcher struct {
	storage *storage
	logger  *zap.Logger
}

func New(options ...Option) (result *Dispatcher, err error) {
	result = &Dispatcher{}
	for _, f := range options {
		f(result)
	}
	if result.storage == nil {
		var s *storage
		s, err = newStorage()
		if err != nil {
			return
		}
		fillStorageWithMockData(s)
		result.storage = s
	}

	return
}
func (t *Dispatcher) Put(i *Item) (err error) {
	err = t.storage.put(i)
	return
}

func (t *Dispatcher) GetAll() (result Items, err error) {
	result, err = t.storage.getAll()
	return
}

func Logger(l *zap.Logger) Option {
	return func(d *Dispatcher) {
		d.logger = l
	}
}
