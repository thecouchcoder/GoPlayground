package transact

import (
	"KeyValueStore/core"
	"fmt"
)

func NewTransactionLogger(tlogType string) (core.TransactionLogger, error) {
	switch tlogType {
	case "file":
		return NewFileTransactionLogger("log.txt")
	case "postgres":
		// not done
		return NewPostgresTransactionLogger(PostgresDBParams{})
	default:
		return nil, fmt.Errorf("unknown transaction logger type: %s", tlogType)
	}
}
