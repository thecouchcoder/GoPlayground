package main

import (
	"fmt"
	"sync"
)

func main() {
	channels := make([]<-chan string, 0)
	for i := 0; i < 3; i++ {
		ch := make(chan string)
		channels = append(channels, ch)

		go func(i int) {
			defer close(ch)
			for j := 0; j < 5; j++ {
				ch <- fmt.Sprintf("Channel %d value %d", i, j)
			}
		}(i)
	}

	dest := fanIn(channels...)

	for val := range dest {
		fmt.Println(val)
	}
}

func fanIn(sources ...<-chan string) <-chan string {
	output := make(chan string)
	var wg sync.WaitGroup

	wg.Add(len(sources))

	for _, source := range sources {
		go func(c <-chan string) {
			defer wg.Done()

			for val := range c {
				output <- val
			}
		}(source)
	}

	go func() {
		fmt.Println("Done!")
		wg.Wait()
		close(output)
	}()

	return output
}
