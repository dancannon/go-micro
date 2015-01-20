package transport

import (
	"errors"
	"fmt"
	"github.com/asim/go-micro/registry"

	log "github.com/Sirupsen/logrus"
)

type Transport interface {
	Requester
	Responder

	Init(reg registry.Registry) error
}

type Requester interface {
	NewRequest(service, endpoint string, payload interface{}) (*Request, error)
}

type Responder interface {
	Register(handler interface{})
	RegisterNamed(name string, handler interface{})
	Start(addr string) (Listener, error)
}

type Listener interface {
	Close() error
}

var (
	transports      map[string]Transport
	ErrNotSupported = errors.New("transport not supported")
)

func Register(name string, r Transport) error {
	if _, exists := transports[name]; exists {
		return fmt.Errorf("Transport '%s' already registered", name)
	}
	transports[name] = r
	log.Debugf("Registering transport '%s'", name)

	return nil
}

func New(name string, reg registry.Registry) (Transport, error) {
	if t, exists := transports[name]; exists {
		log.Debugf("Initializing transport '%s'", name)
		err := t.Init(reg)
		return t, err
	}

	return nil, ErrNotSupported
}
