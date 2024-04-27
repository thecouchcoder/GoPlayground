package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	SpawnProcesses(ctx)
	fmt.Println("Done!")
}

func SpawnProcess(ctx context.Context) {
	dctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	done := make(chan bool)
	go func() {
		RandomProcess(dctx, 1)
		done <- true
	}()
	select {
	case <-dctx.Done():
		fmt.Printf("time out %v", dctx.Err())
	case <-done:
		return
	}
}

func SpawnProcesses(ctx context.Context) {
	var wg sync.WaitGroup
	n := 3
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			dctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
			defer cancel()
			done := make(chan bool)

			go func(dctx context.Context, i int) {
				RandomProcess(dctx, i)
				done <- true
			}(dctx, i)
			select {
			case <-dctx.Done():
				fmt.Printf("Process %d: %v\n", i, dctx.Err())
				wg.Done()
			case <-done:
				wg.Done()
			}
		}(i)
	}

	wg.Wait()

}

func RandomProcess(dctx context.Context, i int) {
	n := rand.Int63n(500)
	fmt.Printf("Process %d: Started for %d\n", i, n)
	time.Sleep(time.Duration(n) * time.Millisecond)
	fmt.Printf("Process %d: Finished\n", i)
}
