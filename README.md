> # retry
>
> Functional mechanism based on context to perform actions repetitively until successful.

[![Build Status](https://travis-ci.org/kamilsk/retry.svg?branch=master)](https://travis-ci.org/kamilsk/retry)
[![Coverage Status](https://coveralls.io/repos/github/kamilsk/retry/badge.svg)](https://coveralls.io/github/kamilsk/retry)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamilsk/retry)](https://goreportcard.com/report/github.com/kamilsk/retry)
[![Exago](https://api.exago.io/badge/rank/github.com/kamilsk/retry)](https://www.exago.io/project/github.com/kamilsk/retry)
[![GoDoc](https://godoc.org/github.com/kamilsk/retry?status.svg)](https://godoc.org/github.com/kamilsk/retry)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](LICENSE)

## Differences from [Rican7/retry](https://github.com/Rican7/retry)

- fixed [bug](https://github.com/Rican7/retry/pull/2) with unexpected infinite loop
  - added transparent mechanism for this purpose as `Infinite` [strategy](strategy/strategy.go#L24-L28)
- added `context` support to cancellation
- added `error` transmission between attempts
  - added `classifier` to handle them (see [classifier](classifier) package)
- added cli tool `retry` which provides functionality to repeat terminal commands (see [cmd/retry](cmd))

## Usage

### HTTP calls with retries and backoff

```go
var response struct {
    ID      int
    Message string
}
client := &http.Client{Timeout: 100 * time.Millisecond}

action := func(attempt uint) error {
    resp, err := client.Do(&http.NewRequest(http.MethodGet, "http://localhost:8080", nil))
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

```go
MustOpen := func() *sql.DB {
	db, err := sql.Open("sqlite", "./sqlite.db")
	if err != nil {
		panic(err)
	}
	return db
}

go func(db *sql.DB, attempt uint, frequency time.Duration) {
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
		if err := retry.Retry(context.Background(), ping, strategies...); err != nil {
			panic(err)
		}
		time.Sleep(frequency)
	}
}(MustOpen(), 10, time.Minute)
```

### CLI tool for command execution repetitively

```bash
$ export PATH=$GOPATH/bin:$PATH
$ go install ./cmd/retry
$ retry -limit=3 -backoff=lin[10ms] -- curl http://unknown.host
curl: (52) Empty reply from server
curl: (52) Empty reply from server
curl: (52) Empty reply from server
```

See more details [here](cmd).

## Installation

```bash
$ egg -version 2.x github.com/kamilsk/retry
```

### Mirror

```bash
$ egg -version 2.x -fix-vanity-url bitbucket.org/kamilsk/retry 
```

> [egg](https://github.com/kamilsk/egg) is an `extended go get`.

### Update

This library is using [SemVer](http://semver.org) for versioning and it is not 
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe.
Therefore, do not use `go get -u` to update it, use [Glide](https://glide.sh) or something similar for this purpose.

## Integration with Docker

```bash
$ make docker-pull
$ make docker-gometalinter ARGS=--deadline=12s
$ make docker-bench ARGS=-benchmem
$ make docker-test ARGS=-v
$ make docker-test-with-coverage ARGS=-v OPEN_BROWSER=true
```

## Feedback

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/kamilsk/retry)
[![@ikamilsk](https://img.shields.io/badge/author-%40ikamilsk-blue.svg)](https://twitter.com/ikamilsk)

## Notes

- tested on Go 1.7 and 1.8, use 1.x version for 1.5 and 1.6
- [research](RESEARCH.md)
