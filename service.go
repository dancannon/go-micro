package micro

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/asim/go-micro/client"
	"github.com/asim/go-micro/registry"
	"github.com/asim/go-micro/server"
	"github.com/asim/go-micro/store"
	"github.com/asim/go-micro/transport"

	_ "github.com/asim/go-micro/registry/consul"
	_ "github.com/asim/go-micro/store/consul"
	_ "github.com/asim/go-micro/transport/rpc"
)

// Service represents a single node in a cluster. A service can be used as both
// a client and server.
type Service struct {
	ID   string
	Name string

	Config Config

	Registry registry.Registry
	Store    store.Store
}

type Config struct {
	Registry     string
	RegistryAddr string

	Store     string
	StoreAddr string

	Transport string
}

func DefaultConfig() Config {
	return Config{
		Registry:  "consul",
		Store:     "consul",
		Transport: "pbrpc",
	}
}

func New(name string, config Config) (*Service, error) {
	if len(name) == 0 {
		name = "go-server"
	}

	id := name + "-" + uuid.NewUUID().String()

	reg, err := registry.New(config.Registry, config.RegistryAddr)
	if err != nil {
		return nil, err
	}
	str, err := store.New(config.Store, config.StoreAddr)
	if err != nil {
		return nil, err
	}

	return &Service{
		ID:   id,
		Name: name,

		Config: config,

		Registry: reg,
		Store:    str,
	}, nil
}

func (s *Service) NewClient() (*client.Client, error) {
	// Get transport
	t, err := transport.New(s.Config.Transport, s.Registry)
	if err != nil {
		return nil, err
	}

	// Get client
	cl, err := client.New(s.ID, s.Name, t, s.Registry)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (s *Service) NewServer() (*server.Server, error) {
	// Get transport
	t, err := transport.New(s.Config.Transport, s.Registry)
	if err != nil {
		return nil, err
	}

	srv, err := server.New(s.ID, s.Name, t, s.Registry)
	if err != nil {
		return nil, err
	}

	return srv, nil
}
