package retry_test

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	. "github.com/kamilsk/retry/v3"
	. "github.com/kamilsk/retry/v3/backoff"
	. "github.com/kamilsk/retry/v3/jitter"
	. "github.com/kamilsk/retry/v3/strategy"
)

func Example() {
	_ = Retry(nil, func(attempt uint) error {
		return nil // Do something that may or may not cause an error
	})
}

func Example_fileOpen() {
	const logFilePath = "/var/log/myapp.log"

	var logFile *os.File

	err := Retry(nil, func(attempt uint) error {
		var err error

		logFile, err = os.Open(logFilePath)

		return err
	})

	if nil != err {
		log.Fatalf("Unable to open file %q with error %q", logFilePath, err)
	}

	_ = logFile.Chdir() // Do something with the file
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

	err := Retry(nil, action, Limit(5), Backoff(Fibonacci(10*time.Millisecond)))
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

	_ = Retry(
		nil,
		action,
		Limit(5),
		BackoffWithJitter(
			BinaryExponential(10*time.Millisecond),
			Deviation(random, 0.5),
		),
	)
}
