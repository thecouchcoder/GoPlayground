package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	throttle := Throttle(DoStuff, 2, 3*time.Millisecond)
	throttle(ctx)
	throttle(ctx)
	throttle(ctx)
	time.Sleep(5 * time.Millisecond)
	throttle(ctx)
	throttle(ctx)
	throttle(ctx)
	throttle(ctx)
	throttle(ctx)
	throttle(ctx)
}

func DoStuff(ctx context.Context) error {
	fmt.Println("Do Stuff")
	return nil
}

func Throttle(call func(ctx context.Context) error, max int, wait time.Duration) func(ctx context.Context) error {
	var once sync.Once
	var tokens = max
	return func(ctx context.Context) error {

		once.Do(func() {
			fmt.Println("Do Once")
			ticker := time.NewTicker(wait)

			go func() {
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						tokens = max
					}
				}
			}()
		})

		if tokens == 0 {
			fmt.Println("Throttled")
			return fmt.Errorf("Throttled")
		}

		tokens--
		return call(ctx)

	}
}
