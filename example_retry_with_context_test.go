// +build go1.7

package retry_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	. "github.com/kamilsk/retry/v3"
	. "github.com/kamilsk/retry/v3/strategy"
)

// This example shows how to use context and retry together.
func Example_retryWithContext() {
	const expected = 3

	attempts := expected
	communication := make(chan error)
	go func() {
		for {
			<-communication
			if attempts == 0 {
				close(communication)
				return
			}
			attempts--
			communication <- errors.New("service unavailable")
		}
	}()

	action := func(uint) error {
		communication <- nil   // ping
		return <-communication // pong
	}
	ctx := WithContext(context.Background(), WithTimeout(time.Second))
	if err := Retry(ctx.Done(), action, Delay(time.Millisecond)); err != nil {
		panic(err)
	}

	fmt.Printf("attempts: %d, expected: %d \n", expected-attempts, expected)
	// Output: attempts: 3, expected: 3
}
