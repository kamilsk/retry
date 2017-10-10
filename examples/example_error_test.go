package examples

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"

	"github.com/kamilsk/retry"
	"github.com/kamilsk/retry/strategy"
)

type Temporary interface {
	Temporary() bool
}

type HttpError struct {
	Code int
}

func (err HttpError) Error() string {
	return fmt.Sprintf("http: status code %d", err.Code)
}

func (err HttpError) Temporary() bool {
	switch err.Code {
	case http.StatusRequestTimeout:
	case http.StatusBadGateway:
	case http.StatusServiceUnavailable:
		return true
	}
	return false
}

// This example shows how to handle errors.
func Example_handleErrors() {
	var repeat int

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		repeat++
		rw.WriteHeader(http.StatusServiceUnavailable)
		_, _ = rw.Write([]byte("Internal Server Error"))
	}))
	defer ts.Close()

	checkNetworkError := func(attempt uint, err error) bool {
		if err, ok := err.(net.Error); ok {
			return err.Timeout() || err.Temporary()
		}
		return true
	}

	checkStatusCode := func(attempt uint, err error) bool {
		if err, ok := err.(Temporary); ok {
			return err.Temporary()
		}
		return true
	}

	action := func(attempt uint) error {
		resp, err := http.Get(ts.URL)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return HttpError{Code: resp.StatusCode}
		}
		return nil
	}

	if err := retry.Retry(nil, action, checkNetworkError, checkStatusCode, strategy.Limit(2)); err != nil {
		fmt.Printf("repeat: %d, err: %q \n", repeat, err)
	}

	// Output: repeat: 2, err: "http: status code 503"
}
