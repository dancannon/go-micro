package registry

import (
	"fmt"
	"os"

	k8s "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
)

type KubernetesRegistry struct {
	Client    *k8s.Client
	Namespace string
}

func init() {
	registry.Register("kubernetes", &KubernetesRegistry{})
}

func (r *KubernetesRegistry) Init(address string) error {
	config := &k8s.Config{
		Host: address,
	}

	client, err := k8s.New(config)
	if err != nil {
		return err
	}

	r.Client = client
	r.Namespace = "default"

	return nil
}

func (c *KubernetesRegistry) Deregister(s registry.Service) error {
	return nil
}

func (c *KubernetesRegistry) Register(s registry.Service) error {
	return nil
}

func (c *KubernetesRegistry) GetService(name string) (registry.Service, error) {
	services, err := c.Client.Services(c.Namespace).List(labels.OneTermEqualSelector("name", name))
	if err != nil {
		return nil, err
	}

	if len(services.Items) == 0 {
		return nil, fmt.Errorf("Service not found")
	}

	ks := &KubernetesService{ServiceName: name}
	for _, item := range services.Items {
		ks.ServiceNodes = append(ks.ServiceNodes, &KubernetesNode{
			NodeAddress: item.Spec.PortalIP,
			NodePort:    item.Spec.Port,
		})
	}

	return ks, nil
}

func (c *KubernetesRegistry) NewService(name string, nodes ...registry.Node) registry.Service {
	var snodes []*KubernetesNode

	for _, node := range nodes {
		if n, ok := node.(*KubernetesNode); ok {
			snodes = append(snodes, n)
		}
	}

	return &KubernetesService{
		ServiceName:  name,
		ServiceNodes: snodes,
	}
}

func (c *KubernetesRegistry) NewNode(id, address string, port int) registry.Node {
	return &KubernetesNode{
		NodeId:      id,
		NodeAddress: address,
		NodePort:    port,
	}
}
