package examples

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/kamilsk/breaker"
	"github.com/kamilsk/retry/v4"
	"github.com/kamilsk/retry/v4/strategy"
)

var server *httptest.Server

func TestMain(m *testing.M) {
	server = httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) { rw.WriteHeader(http.StatusOK) }),
	)
	code := m.Run()
	server.Close()
	os.Exit(code)
}

func ExampleRetryQuickStart() {
	var response *http.Response

	action := func(uint) error {
		var err error
		response, err = http.Get(server.URL)
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
		response, err = http.Get(server.URL)
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
		req, err := http.NewRequest(http.MethodGet, server.URL, nil)
		if err != nil {
			return err
		}
		req = req.WithContext(ctx)
		response, err = http.DefaultClient.Do(req)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	br, ctx := breaker.WithContext(ctx)
	defer func() {
		// they do the same thing
		br.Close()
		cancel()
	}()

	if err := retry.TryContext(ctx, action, strategy.Limit(3)); err != nil {
		log.Fatal(err)
	}

	_, _ = io.Copy(ioutil.Discard, response.Body)
	_ = response.Body.Close()

	fmt.Println(response.Status)
	// Output: 200 OK
}
