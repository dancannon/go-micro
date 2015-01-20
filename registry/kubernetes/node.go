package kubernetes

import (
	"github.com/asim/go-micro/registry"
)

type KubernetesNode struct {
	NodeId      string
	NodeAddress string
	NodePort    int
}

func (c *KubernetesNode) Id() string {
	return c.NodeId
}

func (c *KubernetesNode) Address() string {
	return c.NodeAddress
}

func (c *KubernetesNode) Port() int {
	return c.NodePort
}
