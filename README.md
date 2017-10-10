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

## Code of Conduct

The project team follows [Contributor Covenant v1.4](http://contributor-covenant.org/version/1/4/).
Instances of abusive, harassing or otherwise unacceptable behavior may be reported by contacting
the project team at feedback@octolab.org.

---

## Differences from [Rican7/retry](https://github.com/Rican7/retry)

- Fixed [bug](https://github.com/Rican7/retry/pull/2) with an unexpected infinite loop.
  - Added a clear mechanism for this purpose as the Infinite [strategy](strategy/strategy.go#L24-L28).
- Added `context` support to cancellation.
- Added `error` transmission between attempts.
  - Added `classifier` to handle them (see [classifier](classifier) package).
- Added CLI tool `retry` which provides functionality for repeating terminal commands (see [cmd/retry](cmd)).

## [Usage](examples#Usage)

## Installation

```bash
$ egg github.com/kamilsk/retry
```

### Mirror

```bash
$ egg bitbucket.org/kamilsk/retry
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
$ make check-code-quality
```

### Testing

#### Local

```bash
$ make install-deps
$ make test # or test-with-coverage
$ make bench
```

#### Docker

```bash
$ make docker-pull
$ make complex-tests # or complex-tests-with-coverage
$ make complex-bench
```

## Feedback

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/kamilsk/retry)
[![@ikamilsk](https://img.shields.io/badge/author-%40ikamilsk-blue.svg)](https://twitter.com/ikamilsk)

## Notes

- tested on Go 1.7 and 1.8, use 1.x version for 1.5 and 1.6
- [research](RESEARCH.md)
