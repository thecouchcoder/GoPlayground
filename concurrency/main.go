package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	SendAndReceive()
	fmt.Println("------")
	SendDelayAndReceive()
	fmt.Println("------")
	SendAndReceiveDelay()
	fmt.Println("------")

}

func SendAndReceive() {
	t := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(1)
	ch := make(chan string)
	go func() {
		fmt.Printf("%v: started send and receive go routine\n", int(time.Since(t).Seconds()))
		v := <-ch
		fmt.Printf("%v: received %s\n", int(time.Since(t).Seconds()), v)
	}()

	fmt.Printf("%v: sent\n", int(time.Since(t).Seconds()))
	ch <- "send and receive"
	wg.Done()
}

func SendDelayAndReceive() {
	t := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(1)
	ch := make(chan string)
	go func() {
		fmt.Printf("%v: started send delay and receive go routine\n", int(time.Since(t).Seconds()))
		v := <-ch
		fmt.Printf("%v: received %s\n", int(time.Since(t).Seconds()), v)
		wg.Done()
	}()
	time.Sleep(2 * time.Second)
	fmt.Printf("%v: sent after delay\n", int(time.Since(t).Seconds()))
	ch <- "send delay and receive"

	wg.Wait()
}

func SendAndReceiveDelay() {
	t := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(1)
	ch := make(chan string)
	go func() {
		fmt.Printf("%v: started send delay and receive go routine\n", int(time.Since(t).Seconds()))
		time.Sleep(2 * time.Second)
		v := <-ch
		fmt.Printf("%v: received %s\n", int(time.Since(t).Seconds()), v)
		wg.Done()
	}()

	fmt.Printf("%v: sent expecting receive delay\n", int(time.Since(t).Seconds()))
	ch <- "send  delay and receive"

	wg.Wait()
}
