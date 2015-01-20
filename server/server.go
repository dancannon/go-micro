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

func (s *Server) Start(address string) error {
	l, err := s.transport.Start(address)
	if err != nil {
		return err
	}

	// parse address for host, port
	parts := strings.Split(address, ":")
	host := strings.Join(parts[:len(parts)-1], ":")
	port, _ := strconv.Atoi(parts[len(parts)-1])

	// register service
	node := registry.Node{s.id, host, port}
	service := registry.Service{s.name, []registry.Node{node}}

	log.Debugf("Registering %s", node.Id)
	s.registry.Register(service)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	for {
		select {
		case <-ch:
			log.Debugf("Received signal")
		case <-s.close:
			log.Debugf("Stopping server")
		}
	}

	log.Debugf("Deregistering %s", node.Id)
	s.registry.Deregister(service)

	return l.Close()
}

func (s *Server) Close() error {
	s.close <- struct{}{}
	return nil
}

func (s *Server) Register(handler interface{}) {
	s.transport.Register(handler)
}

func (s *Server) RegisterNamed(name string, handler interface{}) {
	s.transport.RegisterNamed(name, handler)
}
