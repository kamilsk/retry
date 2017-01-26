// +build go1.7

package retrier_test

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/kamilsk/retrier"
	"github.com/kamilsk/retrier/backoff"
	"github.com/kamilsk/retrier/jitter"
	"github.com/kamilsk/retrier/net"
	"github.com/kamilsk/retrier/strategy"
)

func Example() {
	retrier.Retry(func(attempt uint) error {
		return nil // Do something that may or may not cause an error
	})
}

func Example_fileOpen() {
	const logFilePath = "/var/log/myapp.log"

	var logFile *os.File

	err := retrier.Retry(func(attempt uint) error {
		var err error

		logFile, err = os.Open(logFilePath)

		return err
	})

	if nil != err {
		log.Fatalf("Unable to open file %q with error %q", logFilePath, err)
	}

	logFile.Chdir() // Do something with the file
}

func Example_httpGetWithStrategies() {
	var response *http.Response
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))

	action := func(attempt uint) error {
		var err error

		response, err = http.Get(ts.URL)

		if nil == err && nil != response && response.StatusCode > 200 {
			err = fmt.Errorf("failed to fetch (attempt #%d) with status code: %d", attempt, response.StatusCode)
		}

		return err
	}

	err := retrier.Retry(
		action,
		strategy.Limit(5),
		strategy.Backoff(backoff.Fibonacci(10*time.Millisecond)),
	)

	if nil != err {
		log.Fatalf("Failed to fetch repository with error %q", err)
	}
}

func Example_withBackoffJitter() {
	action := func(attempt uint) error {
		return errors.New("something happened")
	}

	seed := time.Now().UnixNano()
	random := rand.New(rand.NewSource(seed))

	retrier.Retry(
		action,
		strategy.Limit(5),
		strategy.BackoffWithJitter(
			backoff.BinaryExponential(10*time.Millisecond),
			jitter.Deviation(random, 0.5),
		),
	)
}

// This example shows usage of strategy.Timeout(time.Duration).
func Example_timeoutToRetry() {
	webServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		time.Sleep(10 * time.Millisecond)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Internal Server Error"))
	}))

	action := func(attempt uint) error {
		resp, err := http.Get(webServer.URL)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return errors.New("bad response")
		}
		return nil
	}

	if err := retrier.Retry(action, strategy.Timeout(30*time.Millisecond), strategy.Delay(20*time.Millisecond)); err != nil {
		// err.Error() == "bad response"
	}
}

// This example shows how to operate on errors.
//
//	type Temporary interface {
//		Temporary() bool
//	}
//
//	type HttpError struct {
//		Code int
//	}
//
//	func (err *HttpError) Error() string {
//		return fmt.Sprintf("http: status code %d", err.Code)
//	}
//
//	func (err *HttpError) Temporary() bool {
//		switch err.Code {
//		case http.StatusRequestTimeout:
//		case http.StatusBadGateway:
//		case http.StatusServiceUnavailable:
//			return true
//		}
//		return false
//	}
func Example_operateOnError() {
	webServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Internal Server Error"))
	}))

	checkStatusCode := func(attempt uint, err error) bool {
		if tempError, ok := err.(Temporary); ok {
			return tempError.Temporary()
		}
		return true
	}

	action := func(attempt uint) error {
		resp, err := http.Get(webServer.URL)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return &HttpError{Code: resp.StatusCode}
		}
		return nil
	}

	if err := retrier.RetryWithError(action, net.CheckNetError(), checkStatusCode); err != nil {
		// this code will not be executed
	}
}

// helpers

type Temporary interface {
	Temporary() bool
}

type HttpError struct {
	Code int
}

func (err *HttpError) Error() string {
	return fmt.Sprintf("http: status code %d", err.Code)
}

func (err *HttpError) Temporary() bool {
	switch err.Code {
	case http.StatusRequestTimeout:
	case http.StatusBadGateway:
	case http.StatusServiceUnavailable:
		return true
	}
	return false
}
