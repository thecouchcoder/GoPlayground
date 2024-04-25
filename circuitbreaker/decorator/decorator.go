package decorator

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Gateway struct{}

func (g Gateway) Execute(req interface{}) (interface{}, error) {
	return nil, fmt.Errorf("bad Request %v", req)
}

type CircuitBreaker interface {
	Execute(req interface{}) (interface{}, error)
}

type State string

const (
	open   State = "open"
	closed State = "closed"
)

var (
	ErrCircuitOpen = errors.New("error the decorator circuit is open")
)

type circuitBreaker struct {
	state        State
	fails        int
	maxThreshold int
	openInterval time.Duration
	openCh       chan interface{}
	wrappee      CircuitBreaker
	mutex        sync.Mutex
}

func NewDecoratorCircuitBreaker(wrapee CircuitBreaker) CircuitBreaker {
	cb := circuitBreaker{
		state:        closed,
		fails:        0,
		maxThreshold: 5,
		openInterval: 50 * time.Millisecond,
		openCh:       make(chan interface{}),
		wrappee:      wrapee,
	}

	go cb.openWatcher()
	return &cb
}

func (cb *circuitBreaker) openWatcher() {
	for range cb.openCh {
		time.Sleep(cb.openInterval)
		cb.mutex.Lock()
		cb.fails = 0
		cb.state = closed
		cb.mutex.Unlock()
	}
}

func (cb *circuitBreaker) Execute(req interface{}) (interface{}, error) {
	err := cb.doPreRequest()
	if err != nil {
		return nil, err
	}

	res, err := cb.wrappee.Execute(req)

	err = cb.doPostRequest(err)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (cb *circuitBreaker) doPreRequest() error {
	if cb.state == open {
		return ErrCircuitOpen
	}

	return nil
}

func (cb *circuitBreaker) doPostRequest(err error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err == nil {
		cb.fails = 0
		cb.state = closed
		return nil
	}

	cb.fails++
	if cb.fails > cb.maxThreshold {
		cb.state = open
		cb.openCh <- 0
	}

	return err
}
