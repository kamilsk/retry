> # retry
>
> A simple, stateless, functional mechanism to perform actions repetitively until successful.
>
> > Extended [Rican7/retry](https://github.com/Rican7/retry) package.

[![Build Status](https://travis-ci.org/kamilsk/retry.svg?branch=master)](https://travis-ci.org/kamilsk/retry)
[![Coverage Status](https://coveralls.io/repos/github/kamilsk/retry/badge.svg)](https://coveralls.io/github/kamilsk/retry)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamilsk/retry)](https://goreportcard.com/report/github.com/kamilsk/retry)
[![GoDoc](https://godoc.org/github.com/kamilsk/retry?status.svg)](https://godoc.org/github.com/kamilsk/retry)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](LICENSE.md)

## What's the differences?

| Rican7/retry        | kamilsk/retry                                          | Description                                    |
|:--------------------|:-------------------------------------------------------|:-----------------------------------------------|
| retry.Retry(action) | retry.Retry(action, strategy.Infinite())               | a clear indication of the infinity of attempts |
| -                   | retry.Retry(action, strategy.Timeout(time.Duration))   | timeout to retry                               |
| -                   | retry.RetryWithError(action, strategies...)            | extended strategy could operate on error       |
| -                   | retry.RetryWithError(action, strategy.CheckNetError()) | handle temporary or timeout network errors     |

## Installation

```bash
$ go get github.com/kamilsk/retrier
```

### Mirror

```bash
$ go get bitbucket.org/kamilsk/retrier
```

## Feedback

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/kamilsk/retry)
[![@ikamilsk](https://img.shields.io/badge/author-%40ikamilsk-blue.svg)](https://twitter.com/ikamilsk)
