package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var ErrorNoSuchKey = errors.New("no such key")

var store = make(map[string]string)

func helloGoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	response := fmt.Sprintf("Hello %s!\n", name)
	w.Write([]byte(response))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/{name}", helloGoHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func Put(key string, value string) error {
	store[key] = value
	return nil
}

func Get(key string) (string, error) {
	value, ok := store[key]
	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func Delete(key string) error {
	delete(store, key)
	return nil
}
