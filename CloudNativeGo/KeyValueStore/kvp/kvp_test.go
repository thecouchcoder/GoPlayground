package kvp

import (
	"errors"
	"testing"
)

func TestPut(t *testing.T) {
	const key = "put-key"
	const value = "put-val"
	defer delete(store, key)

	if _, ok := store[key]; ok {
		t.Error("Found key before put")
	}

	err := Put(key, value)

	if err != nil {
		t.Error(err)
	}

	val, ok := store[key]
	if !ok {
		t.Error("Key not inserted")
	}
	if value != val {
		t.Error("Value doesn't match")
	}
}
func TestGet(t *testing.T) {
	const key = "get-key"
	const value = "get-val"
	defer delete(store, key)

	_, err := Get(key)
	if err == nil {
		t.Error("Expected an error")
	}
	if !errors.Is(err, ErrorNoSuchKey) {
		t.Error("unexpected error", err)
	}

	store[key] = value

	val, err := Get(key)
	if err != nil {
		t.Error("unexpected error", err)
	}
	if val != value {
		t.Error("Wrong value")
	}
}

func TestDelete(t *testing.T) {
	const key = "delete-key"
	const value = "delete-val"
	defer delete(store, key)

	store[key] = value

	if _, contains := store[key]; !contains {
		t.Error("key/value doesn't exist")
	}

	err := Delete(key)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if _, contains := store[key]; contains {
		t.Error("key still exists")
	}
}
