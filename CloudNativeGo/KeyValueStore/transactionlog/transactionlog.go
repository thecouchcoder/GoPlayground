package transactionlog

type TransactionLogger interface {
	LogPut(key string, value string)
	LogDelete(key string)

	Err() <-chan error

	ReadEvents() (<-chan Event, <-chan error)
	Run()
	Close()
}

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}

type EventType byte

const (
	_ EventType = iota
	PUT
	DELETE
)
