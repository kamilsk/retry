> # retry
>
> Functional mechanism based on context to perform actions repetitively until successful.

[![Build Status](https://travis-ci.org/kamilsk/retry.svg?branch=master)](https://travis-ci.org/kamilsk/retry)
[![Coverage Status](https://coveralls.io/repos/github/kamilsk/retry/badge.svg)](https://coveralls.io/github/kamilsk/retry)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamilsk/retry)](https://goreportcard.com/report/github.com/kamilsk/retry)
[![Exago](https://api.exago.io/badge/rank/github.com/kamilsk/retry)](https://www.exago.io/project/github.com/kamilsk/retry)
[![GoDoc](https://godoc.org/github.com/kamilsk/retry?status.svg)](https://godoc.org/github.com/kamilsk/retry)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](LICENSE)

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
$ go get github.com/kamilsk/retry
```

### Mirror

```bash
$ go get bitbucket.org/kamilsk/retry | egg -fix-vanity-url
```

> [egg](https://github.com/kamilsk/egg) is an `extended go get`.

### Update

This library is using [SemVer](http://semver.org) for versioning and it is not [BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe.
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

- tested on Go 1.5, 1.6, 1.7 and 1.8
- [research](RESEARCH.md)
