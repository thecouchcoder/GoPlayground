package frontend

import (
	"KeyValueStore/core"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// Front end port

type FrontEnd interface {
	Start(kv *core.KeyValueStore) error
}

type restFrontEnd struct {
	store *core.KeyValueStore
}

func NewRestFrontEnd() FrontEnd {
	return &restFrontEnd{}
}

func (f *restFrontEnd) helloGoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	response := fmt.Sprintf("Hello %s!\n", name)
	w.Write([]byte(response))
}

func (f *restFrontEnd) getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	result, err := f.store.Get(key)
	if errors.Is(err, core.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Write([]byte(result))
	}
}

func (f *restFrontEnd) putHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = f.store.Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (f *restFrontEnd) deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	err := f.store.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (f *restFrontEnd) Start(kv *core.KeyValueStore) error {
	f.store = kv

	r := mux.NewRouter()

	r.HandleFunc("/{name}", f.helloGoHandler)
	r.HandleFunc("/v1/key/{key}", f.getHandler).Methods("GET")
	r.HandleFunc("/v1/key/{key}", f.putHandler).Methods("PUT")
	r.HandleFunc("/v1/key/{key}", f.deleteHandler).Methods("DELETE")

	return http.ListenAndServe(":8080", r)
}
