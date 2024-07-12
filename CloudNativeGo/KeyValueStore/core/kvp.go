package core

import (
	"errors"
	"sync"
)

var ErrorNoSuchKey = errors.New("no such key")

type distributedMap struct {
	sync.RWMutex
	m map[string]string
}

type KeyValueStore struct {
	sync.RWMutex
	m  map[string]string
	tr TransactionLogger
}

func NewKeyValueStore(tr TransactionLogger) *KeyValueStore {
	return &KeyValueStore{
		m:  make(map[string]string),
		tr: tr,
	}
}

func (store *KeyValueStore) Put(key string, value string) error {
	store.Lock()
	defer store.Unlock()
	store.m[key] = value
	store.tr.LogPut(key, value)
	return nil
}

func (store *KeyValueStore) Get(key string) (string, error) {
	store.RLock()
	defer store.RUnlock()

	value, ok := store.m[key]
	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func (store *KeyValueStore) Delete(key string) error {
	store.Lock()
	defer store.Unlock()

	delete(store.m, key)
	store.tr.LogDelete(key)
	return nil
}

func (store *KeyValueStore) Restore() error {
	var err error

	eventCh, errCh := store.tr.ReadEvents()
	event := Event{}
	more := true
	for more && err == nil {
		select {
		case err = <-errCh:
		case event, more = <-eventCh:
			switch event.EventType {
			case DELETE:
				store.Delete(event.Key)
			case PUT:
				store.Put(event.Key, event.Value)
			}
		}
	}

	store.tr.Run()

	return err
}
