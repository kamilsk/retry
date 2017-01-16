> # retrier
>
> Functional mechanism to perform actions repetitively until successful based on context.

[![Build Status](https://travis-ci.org/kamilsk/retrier.svg?branch=master)](https://travis-ci.org/kamilsk/retrier)
[![Coverage Status](https://coveralls.io/repos/github/kamilsk/retrier/badge.svg)](https://coveralls.io/github/kamilsk/retrier)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamilsk/retrier)](https://goreportcard.com/report/github.com/kamilsk/retrier)
[![GoDoc](https://godoc.org/github.com/kamilsk/retrier?status.svg)](https://godoc.org/github.com/kamilsk/retrier)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](LICENSE.md)

## Usage

### What's the differences?

| Rican7/retry        | kamilsk/retrier                                         | Description                                    |
|:--------------------|:--------------------------------------------------------|:-----------------------------------------------|
| retry.Retry(action) | retrier.Retry(action, strategy.Infinite())              | a clear indication of the infinity of attempts |
| -                   | retrier.Retry(action, strategy.Timeout(time.Duration))  | timeout to retry                               |
| -                   | retrier.RetryWithError(action, strategies...)           | extended strategy could operate on error       |
| -                   | retrier.RetryWithError(action, strategy.CheckNetError() | handle temporary or timeout network errors     |

### More examples are coming soon...

## Installation

```bash
$ go get github.com/kamilsk/retrier
```

### Mirror

```bash
$ go get bitbucket.org/kamilsk/retrier
```

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

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/kamilsk/retrier)
[![@ikamilsk](https://img.shields.io/badge/author-%40ikamilsk-blue.svg)](https://twitter.com/ikamilsk)

## Notes

- tested on Go 1.5, 1.6 and 1.7
- [research](RESEARCH.md)
