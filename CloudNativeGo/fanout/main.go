package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

func main() {
	source := make(chan string)
	dests := fanout(source, 3)

	go func() {
		defer close(source)
		for i := 0; i < 10; i++ {
			source <- fmt.Sprintf("%d", i)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(len(dests))

	for _, ch := range dests {
		go func(ch <-chan string) {
			defer wg.Done()
			for val := range ch {
				fmt.Println(val)
				if !strings.Contains(val, "channel 1") {
					time.Sleep(10 * time.Second)
				}
			}
		}(ch)
	}

	wg.Wait()
}

func fanout(source <-chan string, multiplex int) []<-chan string {
	dests := make([]<-chan string, 0)

	for i := 0; i < multiplex; i++ {
		ch := make(chan string)
		dests = append(dests, ch)

		go func(i int) {
			defer close(ch)
			for val := range source {
				ch <- fmt.Sprintf("Received %v on channel %d", val, i)
			}
		}(i)
	}

	return dests
}
