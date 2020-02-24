package retry_test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/kamilsk/retry/v4"
	"github.com/kamilsk/retry/v4/backoff"
	"github.com/kamilsk/retry/v4/jitter"
	"github.com/kamilsk/retry/v4/strategy"
)

var generator = rand.New(rand.NewSource(0))

func Example() {
	what := func(uint) (err error) { return SendRequest() }

	how := retry.How{
		strategy.Limit(5),
		strategy.BackoffWithJitter(
			backoff.Fibonacci(10*time.Millisecond),
			jitter.NormalDistribution(
				rand.New(rand.NewSource(time.Now().UnixNano())),
				0.25,
			),
		),
		strategy.CheckNetworkError(strategy.Skip),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := retry.Try(ctx, what, how...); err != nil {
		log.Fatal(err)
	}
	fmt.Println("success communication")
	// Output: success communication
}

func SendRequest() error {
	if generator.Intn(5) > 3 {
		return &net.DNSError{Name: "unknown host", IsTemporary: true}
	}
	return nil
}
