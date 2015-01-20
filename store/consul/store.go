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

func (s *ConsulStore) Init(address string) error {
	config := consul.DefaultConfig()
	client, err := consul.NewClient(config)
	if err != nil {
		return err
	}

	s.Client = client

	return nil
}

func (s *ConsulStore) Get(key string) (*store.Item, error) {
	kv, _, err := s.Client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, errors.New("key not found")
	}

	return &store.Item{
		Key:   kv.Key,
		Value: kv.Value,
	}, nil
}

func (s *ConsulStore) Del(key string) error {
	_, err := s.Client.KV().Delete(key, nil)
	return err
}

func (s *ConsulStore) Put(key string, value []byte) error {
	_, err := s.Client.KV().Put(&consul.KVPair{
		Key:   key,
		Value: value,
	}, nil)

	return err
}
