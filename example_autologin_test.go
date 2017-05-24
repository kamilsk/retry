// +build go1.7

package retry_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/kamilsk/retry"
	"github.com/kamilsk/retry/classifier"
	"github.com/kamilsk/retry/strategy"
)

type AuthError struct{}

func (err AuthError) Error() string {
	return "auth needed"
}

func Callback(c classifier.Classifier, f func()) strategy.Strategy {
	return func(attempt uint, err error) bool {
		action := c.Classify(err)
		if action == classifier.Retry {
			f()
		}
		// skip to other strategies if not fail
		return action != classifier.Fail
	}
}

var AuthChecker classifier.FunctionalClassifier = func(err error) classifier.Action {
	if err == nil {
		return classifier.Succeed
	}

	if _, is := err.(AuthError); is {
		return classifier.Retry
	}

	return classifier.Unknown
}

// This example shows how to extend the library to solve a problem with authentication.
func Example_autologin() {
	tokens := map[string]time.Time{}

	mux := http.NewServeMux()
	mux.HandleFunc("/login", func(rw http.ResponseWriter, req *http.Request) {
		tokens[req.Header.Get("token")] = time.Now().Add(10 * time.Millisecond)
	})
	mux.HandleFunc("/api", func(rw http.ResponseWriter, req *http.Request) {
		if ttl, found := tokens[req.Header.Get("token")]; !found || time.Now().After(ttl) {
			rw.Header().Add("WWW-Authenticate", `Bearer realm="/login"`)
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		rw.Write([]byte("success"))
	})
	ts := httptest.NewServer(mux)

	action := func(attempt uint) error {
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/api", nil)
		if err != nil {
			return err
		}

		req.Header.Add("token", "secret")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusUnauthorized {
				return AuthError{}
			}
			return HttpError{Code: resp.StatusCode}
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil
		}

		fmt.Println(string(body))
		// Output: success

		return nil
	}
	ctx := context.Background()
	login := func() {
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/login", nil)
		if err != nil {
			return
		}

		req.Header.Add("token", "secret")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
	}

	if err := retry.Retry(ctx, action, Callback(AuthChecker, login), strategy.Limit(10)); err != nil {
		fmt.Println(err)
	}
}
