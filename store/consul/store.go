package store

import (
	"errors"

	consul "github.com/armon/consul-api"

	"github.com/asim/go-micro/store"
)

type ConsulStore struct {
	Client *consul.Client
}

func init() {
	store.Register("consul", &ConsulStore{})
}

func (c *ConsulStore) Get(key string) (store.Item, error) {
	kv, _, err := c.Client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, errors.New("key not found")
	}

	return store.Item{
		key:   kv.Key,
		value: kv.Value,
	}, nil
}

func (c *ConsulStore) Del(key string) error {
	_, err := c.Client.KV().Delete(key, nil)
	return err
}

func (c *ConsulStore) Put(key string, value []byte) error {
	_, err := c.Client.KV().Put(&consul.KVPair{
		Key:   key,
		Value: value,
	}, nil)

	return err
}
