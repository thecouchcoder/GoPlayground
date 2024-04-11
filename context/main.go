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
	dctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	RandomProcess(dctx, 1)
}

func SpawnProcesses(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		dctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		go func(dctx context.Context, i int) {
			defer cancel()
			RandomProcess(dctx, i)
			wg.Done()
		}(dctx, i)
	}

	wg.Wait()
}

func RandomProcess(dctx context.Context, i int) {
	n := rand.Int63n(11)
	fmt.Printf("Process %d: Started for %d\n", i, n)
	for {
		select {
		case <-dctx.Done():
			fmt.Printf("Process %d: Timed out\n", i)
			return
		default:
			if n == 0 {
				fmt.Printf("Process %d: Finished\n", i)
				return

			}
			time.Sleep(1 * time.Second)
			fmt.Printf("Process %d: Running for %d more\n", i, n)
			n--
		}
	}
}
