// +build example

package examples

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kamilsk/breaker"
	"github.com/kamilsk/retry/v4"
	"github.com/kamilsk/retry/v4/strategy"
)

func ExampleRetryQuickStart() {
	var response *http.Response

	action := func(uint) error {
		var err error
		response, err = http.Get("https://github.com/kamilsk/retry")
		return err
	}

	if err := retry.Retry(breaker.BreakByTimeout(time.Minute), action, strategy.Limit(3)); err != nil {
		log.Fatal(err)
	}

	_, _ = io.Copy(ioutil.Discard, response.Body)
	_ = response.Body.Close()

	fmt.Println(response.Status)
	// Output: 200 OK
}

func ExampleTryQuickStart() {
	var response *http.Response

	action := func(uint) error {
		var err error
		response, err = http.Get("https://github.com/kamilsk/retry")
		return err
	}
	interrupter := breaker.MultiplexTwo(
		breaker.BreakByTimeout(time.Minute),
		breaker.BreakBySignal(os.Interrupt),
	)
	defer interrupter.Close()

	if err := retry.Try(interrupter, action, strategy.Limit(3)); err != nil {
		log.Fatal(err)
	}

	_, _ = io.Copy(ioutil.Discard, response.Body)
	_ = response.Body.Close()

	fmt.Println(response.Status)
	// Output: 200 OK
}

func ExampleTryContextQuickStart() {
	var response *http.Response

	action := func(ctx context.Context, _ uint) error {
		req, err := http.NewRequest(http.MethodGet, "https://github.com/kamilsk/retry", nil)
		if err != nil {
			return err
		}
		req = req.WithContext(ctx)
		response, err = http.DefaultClient.Do(req)
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	br, ctx := breaker.WithContext(ctx)
	defer br.Close()

	if err := retry.TryContext(ctx, action, strategy.Limit(3)); err != nil {
		log.Fatal(err)
	}

	_, _ = io.Copy(ioutil.Discard, response.Body)
	_ = response.Body.Close()

	fmt.Println(response.Status)
	// Output: 200 OK
}
