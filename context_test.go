// +build go1.7

package retry_test

import (
	"context"
	"testing"
	"time"

	"github.com/kamilsk/retry"
)

func TestWithContext(t *testing.T) {
	sleep := 100 * time.Millisecond
	ctx := retry.WithContext(context.TODO(), retry.WithTimeout(sleep))

	start := time.Now()
	<-ctx.Done()
	end := time.Now()

	if expected, obtained := sleep, end.Sub(start); expected > obtained {
		t.Errorf("an unexpected sleep time. expected: %v; obtained: %v", expected, obtained)
	}
}
