package kvp

import (
	"errors"
	"sync"
)

var ErrorNoSuchKey = errors.New("no such key")

type distributedMap struct {
	sync.RWMutex
	m map[string]string
}

var store = distributedMap{
	m: make(map[string]string),
}

func Put(key string, value string) error {
	store.Lock()
	defer store.Unlock()
	store.m[key] = value
	return nil
}

func Get(key string) (string, error) {
	store.RLock()
	defer store.RUnlock()

	value, ok := store.m[key]
	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func Delete(key string) error {
	store.Lock()
	defer store.Unlock()

	delete(store.m, key)
	return nil
}
