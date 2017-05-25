> # retry
>
> Functional mechanism based on context to perform actions repetitively until successful.

[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#utilities)
[![Build Status](https://travis-ci.org/kamilsk/retry.svg?branch=master)](https://travis-ci.org/kamilsk/retry)
[![Coverage Status](https://coveralls.io/repos/github/kamilsk/retry/badge.svg)](https://coveralls.io/github/kamilsk/retry)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamilsk/retry)](https://goreportcard.com/report/github.com/kamilsk/retry)
[![Exago](https://api.exago.io/badge/rank/github.com/kamilsk/retry)](https://www.exago.io/project/github.com/kamilsk/retry)
[![GoDoc](https://godoc.org/github.com/kamilsk/retry?status.svg)](https://godoc.org/github.com/kamilsk/retry)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](LICENSE)

## Differences from [Rican7/retry](https://github.com/Rican7/retry)

- Fixed [bug](https://github.com/Rican7/retry/pull/2) with an unexpected infinite loop.
  - Added a clear mechanism for this purpose as the Infinite [strategy](strategy/strategy.go#L24-L28).
- Added `context` support to cancellation.
- Added `error` transmission between attempts.
  - Added `classifier` to handle them (see [classifier](classifier) package).
- Added CLI tool `retry` which provides functionality for repeating terminal commands (see [cmd/retry](cmd)).

## Usage

### HTTP calls with retries and backoff

This example shows how to repeat http calls.

```go
var response struct {
    ID      int
    Message string
}
client := &http.Client{Timeout: 100 * time.Millisecond}

action := func(attempt uint) error {
    resp, err := client.Do(&http.NewRequest(http.MethodGet, "http://some.json.api", nil))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    return json.Unmarshal(data, &response)
}

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := retry.Retry(ctx, action, strategy.Backoff(backoff.Exponential(100*time.Millisecond, math.Pi))); err != nil {
    // handle error
}
// handle response
```

### Database connection restore

This example shows how to use the library to restore database connection.

```go
MustOpen := func() *sql.DB {
	db, err := sql.Open("sqlite", "./sqlite.db")
	if err != nil {
		panic(err)
	}
	return db
}

go func(db *sql.DB, ctx context.Context, shutdown chan<- struct{}, attempt uint, frequency time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			shutdown <- struct{}{}
		}
	}()

	ping := func(attempt uint) error {
		return db.Ping()
	}
	strategies := []strategy.Strategy{
		strategy.Limit(attempt),
		strategy.BackoffWithJitter(
			backoff.Incremental(100*time.Millisecond, time.Second),
			jitter.NormalDistribution(rand.New(rand.NewSource(time.Now().UnixNano())), 2.0),
		),
	}

	for {
		if err := retry.Retry(ctx, ping, strategies...); err != nil {
			panic(err)
		}
		time.Sleep(frequency)
	}
}(MustOpen(), context.Background(), shutdown, 10, time.Minute)
```

### Autologin

This example shows how to extend the library to solve a problem with authentication using classifier.

```go
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

client := &http.Client{Timeout: 100 * time.Millisecond}

action := func(attempt uint) error {
	req, err := http.NewRequest(http.MethodGet, "http://some.api/get", nil)
	if err != nil {
		return err
	}

	req.Header.Add("token", "secret")
	resp, err := client.Do(req)
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

	// handle response

	return nil
}
ctx := context.Background()
login := func() {
	req, err := http.NewRequest(http.MethodGet, "http://some.api/login", nil)
	if err != nil {
		return
	}

	req.Header.Add("token", "secret")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// handle response
}

if err := retry.Retry(ctx, action, Callback(AuthChecker, login), strategy.Limit(10)); err != nil {
	// handle error
}
```

### CLI tool for command execution repetitively

```bash
$ retry -limit=3 -backoff=lin{10ms} -- /bin/sh -c 'echo "trying..."; exit 1'
trying...
trying...
trying...
[ERROR] error occurred: "exit status 1"
```

See more details [here](cmd).

## Installation

```bash
$ go get github.com/kamilsk/retry
```

### Mirror

```bash
$ egg -fix-vanity-url bitbucket.org/kamilsk/retry 
```

> [egg](https://github.com/kamilsk/egg) is an `extended go get`.

### Update

This library is using [SemVer](http://semver.org) for versioning and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe.
Therefore, do not use `go get -u` to update it, use [Glide](https://glide.sh) or something similar for this purpose.

## Contributing workflow

### Code quality checking

```bash
$ make docker-pull-tools
$ make docker-check
```

### Testing

#### Local

```bash
$ make install-deps
$ make test-with-coverage
```

#### Docker

```bash
$ make docker-pull
$ make docker-test-with-coverage
```

## Feedback

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/kamilsk/retry)
[![@ikamilsk](https://img.shields.io/badge/author-%40ikamilsk-blue.svg)](https://twitter.com/ikamilsk)

## Notes

- tested on Go 1.5, 1.6, 1.7 and 1.8
- [research](RESEARCH.md)
