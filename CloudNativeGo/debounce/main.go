package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	debounce := Debounce(DoStuff, 100*time.Second)
	fmt.Println(debounce())
	fmt.Println(debounce())
	fmt.Println(debounce())
}

func DoStuff() int {
	fmt.Println("Do Stuff")
	return 42
}

func Debounce(call func() int, wait time.Duration) func() int {
	var next time.Time
	mutex := sync.Mutex{}
	var result int

	return func() int {

		mutex.Lock()

		defer func() {
			next = time.Now().Add(wait)
			mutex.Unlock()
		}()

		if !time.Now().Before(next) {
			callResult := call()
			result = callResult
		}

		return result
	}
}
