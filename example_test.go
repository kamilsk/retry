package retry_test

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/kamilsk/retry/v5"
	"github.com/kamilsk/retry/v5/backoff"
	"github.com/kamilsk/retry/v5/exp"
	"github.com/kamilsk/retry/v5/jitter"
	"github.com/kamilsk/retry/v5/strategy"
)

var generator = rand.New(rand.NewSource(0))

func Example() {
	what := SendRequest

	how := retry.How{
		strategy.Limit(5),
		strategy.BackoffWithJitter(
			backoff.Fibonacci(10*time.Millisecond),
			jitter.NormalDistribution(
				rand.New(rand.NewSource(time.Now().UnixNano())),
				0.25,
			),
		),

		// experimental
		exp.CheckError(
			exp.NetworkError(exp.Skip),
			DatabaseError(),
		),
	}

	breaker, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := retry.Do(breaker, what, how...); err != nil {
		panic(err)
	}
	fmt.Println("success communication")
	// Output: success communication
}

func SendRequest(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if generator.Intn(5) > 3 {
		return &net.DNSError{Name: "unknown host", IsTemporary: true}
	}
	return nil
}

func DatabaseError() func(error) bool {
	blacklist := []error{sql.ErrNoRows, sql.ErrConnDone, sql.ErrTxDone}
	return func(err error) bool {
		for _, preset := range blacklist {
			if err == preset {
				return false
			}
		}
		return true
	}
}
