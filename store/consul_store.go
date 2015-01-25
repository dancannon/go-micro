package store

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
)

type ConsulStore struct {
	Client *consul.Client
}

func (c *ConsulStore) Get(key string) (Item, error) {
	kv, _, err := c.Client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, fmt.Errorf("key %s not found", key)
	}

	return &ConsulItem{
		key:   kv.Key,
		value: kv.Value,
	}, nil
}

func (c *ConsulStore) Del(key string) error {
	_, err := c.Client.KV().Delete(key, nil)
	return err
}

func (c *ConsulStore) Put(item Item) error {
	_, err := c.Client.KV().Put(&consul.KVPair{
		Key:   item.Key(),
		Value: item.Value(),
	}, nil)

	return err
}

func (c *ConsulStore) NewItem(key string, value []byte) Item {
	return &ConsulItem{
		key:   key,
		value: value,
	}
}

func NewConsulStore() Store {
	client, _ := consul.NewClient(consul.DefaultConfig())

	return &ConsulStore{
		Client: client,
	}
}
