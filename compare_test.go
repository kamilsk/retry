package retry_test

import (
	"context"
	"fmt"
	"time"

	"github.com/kamilsk/retry/v5"
	"github.com/kamilsk/retry/v5/strategy"
)

// The example shows the difference between Do and DoAsync
// - if an action doesn't support the interrupt mechanism
// - if a strategy doesn't support the interrupt mechanism
func ExampleDo_badCases() {
	var (
		realTime = 100 * time.Millisecond
		needTime = 5 * time.Millisecond
	)
	{
		badAction := func() error {
			time.Sleep(realTime)
			return nil
		}
		now := time.Now()
		breaker, cancel := context.WithTimeout(context.Background(), needTime)

		Silent(retry.Do(breaker, badAction))
		if time.Since(now) < realTime {
			fmt.Println("unexpected waiting time")
		}
		cancel()
	}
	{
		badStrategy := func(strategy.Breaker, uint, error) bool {
			time.Sleep(realTime)
			return true
		}
		now := time.Now()
		breaker, cancel := context.WithTimeout(context.Background(), needTime)

		Silent(retry.Do(breaker, func() error { return nil }, badStrategy))
		if time.Since(now) < realTime {
			fmt.Println("unexpected waiting time")
		}
		cancel()
	}

	fmt.Println("done")
	// Output: done
}

// The example shows the difference between Do and DoAsync
// - if an action doesn't support the interrupt mechanism
// - if a strategy doesn't support the interrupt mechanism
func ExampleDoAsync_guarantees() {
	var (
		sleepTime  = 100 * time.Millisecond
		needTime   = 5 * time.Millisecond
		inaccuracy = time.Millisecond
	)
	{
		badAction := func() error {
			time.Sleep(sleepTime)
			return nil
		}
		now := time.Now()
		breaker, cancel := context.WithTimeout(context.Background(), needTime)

		Silent(retry.DoAsync(breaker, badAction))
		if time.Since(now)-needTime > time.Millisecond+inaccuracy {
			fmt.Println("unexpected waiting time")
		}
		cancel()
	}
	{
		badStrategy := func(strategy.Breaker, uint, error) bool {
			time.Sleep(sleepTime)
			return true
		}
		now := time.Now()
		breaker, cancel := context.WithTimeout(context.Background(), needTime)

		Silent(retry.DoAsync(breaker, func() error { return nil }, badStrategy))
		if time.Since(now)-needTime > time.Millisecond+inaccuracy {
			fmt.Println("unexpected waiting time")
		}
		cancel()
	}

	fmt.Println("done")
	// Output: done
}

func Silent(error) {}
