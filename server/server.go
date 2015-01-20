package server

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	log "github.com/cihub/seelog"

	"github.com/asim/go-micro/registry"
	"github.com/asim/go-micro/transport"
)

type Server struct {
	id, name  string
	transport transport.Transport
	registry  registry.Registry

	close chan struct{}
}

func New(id, name string, t transport.Transport, r registry.Registry) (*Server, error) {
	return &Server{
		id:        id,
		name:      name,
		transport: t,
		registry:  r,
	}, nil
}

func (s *Server) ListenAndServe(address string) error {
	err := s.transport.StartListening(address)
	if err != nil {
		return err
	}

	// parse address for host, port
	parts := strings.Split(address, ":")
	host := strings.Join(parts[:len(parts)-1], ":")
	port, _ := strconv.Atoi(parts[len(parts)-1])

	// register service
	node := &registry.Node{
		ID:      s.id,
		Address: host,
		Port:    port,
	}
	service := &registry.Service{
		Name:  s.name,
		Nodes: []*registry.Node{node},
	}

	log.Debugf("Registering %s", node.ID)
	err = s.registry.Register(service)
	if err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	log.Debugf("Received signal %s", <-ch)

	log.Debugf("Deregistering %s", node.ID)
	err = s.registry.Deregister(service)
	if err != nil {
		return err
	}

	return s.transport.StopListening()
}

func (s *Server) Close() error {
	return s.transport.StopListening()
}

func (s *Server) Register(handler interface{}) {
	s.transport.Register(handler)
}

func (s *Server) RegisterNamed(name string, handler interface{}) {
	s.transport.RegisterNamed(name, handler)
}
