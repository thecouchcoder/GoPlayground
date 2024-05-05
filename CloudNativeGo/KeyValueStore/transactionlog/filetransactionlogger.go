package transactionlog

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type FileTransactionLogger struct {
	file         *os.File
	lastSequence uint64
	events       chan<- Event // write only channel
	errors       <-chan error // read only channel
}

func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	return &FileTransactionLogger{file: file}, nil
}

func (l *FileTransactionLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		for e := range events {
			l.lastSequence++

			// Write to file
			_, err := fmt.Fprintf(
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan Event)
	outErrors := make(chan error, 1)

	go func() {
		var e Event

		defer close(outEvent)
		defer close(outErrors)

		for scanner.Scan() {
			line := scanner.Text()
			if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s", &e.Sequence, &e.EventType, &e.Key, &e.Value); err != nil {
				outErrors <- err
				return
			}

			if l.lastSequence >= e.Sequence {
				outErrors <- errors.New("sequence number not in sync")
				return
			}

			l.lastSequence = e.Sequence
			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outErrors <- err
			return
		}
	}()

	return outEvent, outErrors
}

func (l *FileTransactionLogger) Close() {
	l.file.Close()
}

func (l *FileTransactionLogger) LogPut(key string, value string) {
	ev := Event{
		EventType: PUT,
		Key:       key,
		Value:     value,
	}
	l.events <- ev
}
func (l *FileTransactionLogger) LogDelete(key string) {
	ev := Event{
		EventType: DELETE,
		Key:       key,
	}
	l.events <- ev
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}
