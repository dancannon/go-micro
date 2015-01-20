package registry

import (
	"errors"

	consul "github.com/armon/consul-api"

	"github.com/asim/go-micro/registry"
)

var (
	ConsulCheckTTL = "30s"
)

type ConsulRegistry struct {
	Client *consul.Client
}

func init() {
	registry.Register("consul", new(ConsulRegistry))
}

func (r *ConsulRegistry) Init(address string) error {
	config := consul.DefaultConfig()
	if address != "" {
		config.Address = address
	}

	client, err := consul.NewClient(config)
	if err != nil {
		return err
	}

	r.Client = client

	return nil
}

func (r *ConsulRegistry) Deregister(s *registry.Service) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one node")
	}

	node := s.Nodes[0]

	_, err := r.Client.Catalog().Deregister(&consul.CatalogDeregistration{
		Node:      node.ID,
		Address:   node.Address,
		ServiceID: node.ID,
	}, nil)

	return err
}

func (r *ConsulRegistry) Register(s *registry.Service) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one node")
	}

	node := s.Nodes[0]

	_, err := r.Client.Catalog().Register(&consul.CatalogRegistration{
		Node:    node.ID,
		Address: node.Address,
		Service: &consul.AgentService{
			Service: s.Name,
			ID:      node.ID,
			Port:    node.Port,
		},
	}, nil)

	return err
}

func (r *ConsulRegistry) GetService(name string) (*registry.Service, error) {
	rsp, _, err := r.Client.Catalog().Service(name, "", nil)
	if err != nil {
		return nil, err
	}

	cs := &registry.Service{}

	for _, s := range rsp {
		if s.ServiceName != name {
			continue
		}

		cs.Name = s.ServiceName
		cs.Nodes = append(cs.Nodes, &registry.Node{
			ID:      s.ServiceID,
			Address: s.Address,
			Port:    s.ServicePort,
		})
	}

	return cs, nil
}

func (r *ConsulRegistry) NewService(name string, nodes ...*registry.Node) *registry.Service {
	return &registry.Service{
		Name:  name,
		Nodes: nodes,
	}
}

func (r *ConsulRegistry) NewNode(id, address string, port int) *registry.Node {
	return &registry.Node{
		ID:      id,
		Address: address,
		Port:    port,
	}
}
