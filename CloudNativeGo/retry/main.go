package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
)

var count = 0

func main() {
	ctx := context.Background()

	fmt.Println("Start retry scenario")
	r := Retry(DoStuff, 5, 1*time.Millisecond)
	r(ctx)

	fmt.Println()
	fmt.Println()
	fmt.Println("Start always fail scenario")
	fail := Retry(AlwaysFail, 5, 1*time.Millisecond)
	fail(ctx)
}

func DoStuff(ctx context.Context) (int, error) {
	count++

	if count <= 3 {
		fmt.Println("Failure")
		return 0, fmt.Errorf("Error")
	}

	fmt.Println("Do Stuff")
	return 42, nil
}

func AlwaysFail(ctx context.Context) (int, error) {
	fmt.Println("Failure")
	return 0, fmt.Errorf("Error")
}

func Retry(call func(ctx context.Context) (int, error), retries int, wait time.Duration) func(context.Context) (int, error) {
	r := 0
	return func(ctx context.Context) (int, error) {
		for {
			result, err := call(ctx)
			r++
			if err == nil || r > retries {
				return result, err
			}

			fmt.Printf("Retry %d\n", r)

			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			case <-time.After(wait):
			}

		}
	}
}
