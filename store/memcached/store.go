package store

import (
	"errors"
	"os"

	mc "github.com/bradfitz/gomemcache/memcache"

	"github.com/asim/go-micro/store"
)

type MemcacheStore struct {
	Client *mc.Client
}

func init() {
	store.Register("memcached", &MemcacheStore{})
}

func (r *MemcacheStore) Init(address string) error {
	if address == "" {
		address = "127.0.0.1:11211"
	}

	r.Client = mc.New(address)

	return nil
}

func (r *MemcacheStore) Get(key string) (store.Item, error) {
	kv, err := r.Client.Get(key)
	if err != nil && err == mc.ErrCacheMiss {
		return nil, errors.New("key not found")
	} else if err != nil {
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

func (r *MemcacheStore) Del(key string) error {
	return r.Client.Delete(key)
}

func (r *MemcacheStore) Put(key string, value []byte) error {
	return r.Client.Set(&mc.Item{
		Key:   key,
		Value: value,
	})
}
