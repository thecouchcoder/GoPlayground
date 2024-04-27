package main

import (
	"circuitbreaker/decorator"
	"errors"
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Startup")
	cb := NewCircuitBreaker()
	for i := range 100 {
		callService(cb, i)
		time.Sleep(1 * time.Millisecond)
	}

	decorator := decorator.NewDecoratorCircuitBreaker(decorator.Gateway{})
	for i := range 100 {
		callDecoratorService(decorator, i)
		time.Sleep(1 * time.Millisecond)
	}
}

func callDecoratorService(cb decorator.CircuitBreaker, req int) error {
	_, err := cb.Execute(req)

	if err != nil {
		fmt.Printf("%v: got error %v\n", req, err)
	} else {
		fmt.Println("Success")
	}

	return err
}

func callService(cb CircuitBreaker, req int) {
	_, err := cb.Execute(func() (interface{}, error) {
		return externalCall(req)
	})

	if err != nil {
		fmt.Printf("%v: got error %v\n", req, err)
	} else {
		fmt.Println("Success")
	}
}

func externalCall(req int) (interface{}, error) {
	return nil, fmt.Errorf("bad Request %v", req)
}

type CircuitBreaker interface {
	Execute(func() (interface{}, error)) (interface{}, error)
}

type State string

const (
	open   State = "open"
	closed State = "closed"
)

var (
	ErrCircuitOpen = errors.New("error the circuit is open")
)

type circuitBreaker struct {
	state        State
	fails        int
	maxThreshold int
	openInterval time.Duration
	openCh       chan interface{}
	mutex        sync.Mutex
}

func NewCircuitBreaker() CircuitBreaker {
	cb := circuitBreaker{
		state:        closed,
		fails:        0,
		maxThreshold: 5,
		openInterval: 50 * time.Millisecond,
		openCh:       make(chan interface{}),
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

func (cb *circuitBreaker) Execute(req func() (interface{}, error)) (interface{}, error) {
	err := cb.doPreRequest()
	if err != nil {
		return nil, err
	}

	res, err := req()

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
