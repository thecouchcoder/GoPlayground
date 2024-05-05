package main

import (
	"KeyValueStore/kvp"
	"KeyValueStore/transactionlog"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var logger transactionlog.TransactionLogger

func main() {
	defer logger.Close()
	r := GetRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}

func GetRouter() http.Handler {
	err := initializeTransactionLogger()
	if err != nil {
		panic("can't initialize")
	}
	r := mux.NewRouter()

	r.HandleFunc("/{name}", helloGoHandler)
	r.HandleFunc("/v1/key/{key}", getHandler).Methods("GET")
	r.HandleFunc("/v1/key/{key}", putHandler).Methods("PUT")
	r.HandleFunc("/v1/key/{key}", deleteHandler).Methods("DELETE")

	return r
}

func initializeTransactionLogger() error {
	var err error
	logger, err = transactionlog.NewFileTransactionLogger("log.txt")
	if err != nil {
		return err
	}

	eventCh, errCh := logger.ReadEvents()
	event := transactionlog.Event{}
	more := true
	for more && err == nil {
		select {
		case err = <-errCh:
		case event, more = <-eventCh:
			switch event.EventType {
			case transactionlog.DELETE:
				kvp.Delete(event.Key)
			case transactionlog.PUT:
				kvp.Put(event.Key, event.Value)
			}
		}
	}

	logger.Run()

	return err
}

func helloGoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	response := fmt.Sprintf("Hello %s!\n", name)
	w.Write([]byte(response))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	result, err := kvp.Get(key)
	if errors.Is(err, kvp.ErrorNoSuchKey) {
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

func putHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = kvp.Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.LogPut(key, string(value))

	w.WriteHeader(http.StatusCreated)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	err := kvp.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.LogDelete(key)
}
