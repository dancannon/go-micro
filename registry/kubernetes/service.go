package kubernetes

import (
	"github.com/asim/go-micro/registry"
)

type KubernetesService struct {
	ServiceName  string
	ServiceNodes []*KubernetesNode
}

func (c *KubernetesService) Name() string {
	return c.ServiceName
}

func (c *KubernetesService) Nodes() []registry.Node {
	var nodes []registry.Node

	for _, node := range c.ServiceNodes {
		nodes = append(nodes, node)
	}

	return nodes
}
