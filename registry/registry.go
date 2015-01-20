package registry

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

type Node struct {
	ID      string
	Address string
	Port    int
}

type Service struct {
	Name  string
	Nodes []*Node
}

type Registry interface {
	Init(address string) error
	Register(*Service) error
	Deregister(*Service) error
	GetService(string) (*Service, error)
	NewService(string, ...*Node) *Service
	NewNode(string, string, int) *Node
}

var (
	registries      map[string]Registry
	ErrNotSupported = errors.New("registry not supported")
)

func init() {
	registries = make(map[string]Registry)
}

func Register(name string, r Registry) error {
	if _, exists := registries[name]; exists {
		return fmt.Errorf("Registry '%s' already registered", name)
	}
	registries[name] = r
	log.Debugf("Registering registry '%s'", name)

	return nil
}

func New(name, address string) (Registry, error) {
	if r, exists := registries[name]; exists {
		log.Debugf("Initializing registry '%s'", name)
		err := r.Init(address)
		return r, err
	}

	return nil, ErrNotSupported
}
