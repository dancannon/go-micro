package store

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

type Store interface {
	Init(address string) error

	Get(key string) (*Item, error)
	Del(key string) error
	Put(key string, value []byte) error
}

var (
	stores          map[string]Store
	ErrNotSupported = errors.New("store not supported")
)

func init() {
	stores = make(map[string]Store)
}

func Register(name string, r Store) error {
	if _, exists := stores[name]; exists {
		return fmt.Errorf("Store '%s' already registered", name)
	}
	stores[name] = r
	log.Debugf("Registering store '%s'", name)

	return nil
}

func New(name, address string) (Store, error) {
	if s, exists := stores[name]; exists {
		log.Debugf("Initializing store '%s'", name)
		err := s.Init(address)
		return s, err
	}

	return nil, ErrNotSupported
}
