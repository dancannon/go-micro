package client

import (
	"github.com/asim/go-micro/registry"
	"github.com/asim/go-micro/transport"
)

type Client struct {
	id, name  string
	transport transport.Transport
	registry  registry.Registry
}

func New(id, name string, t transport.Transport, r registry.Registry) (*Client, error) {
	return &Client{
		id:        id,
		name:      name,
		transport: t,
		registry:  r,
	}, nil
}

func (c *Client) NewRequest(
	service, endpoint string, payload interface{},
) (*transport.Request, error) {
	return c.transport.NewRequest(service, endpoint, payload)
}
