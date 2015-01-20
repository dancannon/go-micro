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

func (m *MemcacheStore) Get(key string) (store.Item, error) {
	kv, err := m.Client.Get(key)
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

func (m *MemcacheStore) Del(key string) error {
	return m.Client.Delete(key)
}

func (m *MemcacheStore) Put(key string, value []byte) error {
	return m.Client.Set(&mc.Item{
		Key:   key,
		Value: value,
	})
}
