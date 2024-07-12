package main

import (
	"KeyValueStore/core"
	"KeyValueStore/frontend"
	"KeyValueStore/transact"
	"os"
)

func main() {
	tl, err := transact.NewTransactionLogger(os.Getenv("TLOG_TYPE"))
	if err != nil {
		panic(err)
	}
	store := core.NewKeyValueStore(tl)
	store.Restore()

	fe := frontend.NewRestFrontEnd()
	err = fe.Start(store)
	if err != nil {
		panic(err)
	}
}
