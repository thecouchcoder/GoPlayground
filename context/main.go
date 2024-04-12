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
	err := SpawnProcess(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done!")
}

func SpawnProcess(ctx context.Context) error {
	dctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	done := make(chan bool)
	go func() {
		RandomProcess(dctx, 1)
		done <- true
	}()
	select {
	case <-dctx.Done():
		return fmt.Errorf("time out %v", dctx.Err())
	case <-done:
		return nil
	}

}

func SpawnProcesses(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		dctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
		go func(dctx context.Context, i int) {
			defer cancel()
			RandomProcess(dctx, i)
			wg.Done()
		}(dctx, i)
	}

	wg.Wait()
}

func RandomProcess(dctx context.Context, i int) {
	n := rand.Int63n(500)
	fmt.Printf("Process %d: Started for %d\n", i, n)
	time.Sleep(time.Duration(n) * time.Millisecond)
	fmt.Printf("Process %d: Finished\n", i)
}
