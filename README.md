> # retry
>
> Functional mechanism based on context to perform actions repetitively until successful.

[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#utilities)
[![Patreon](https://img.shields.io/badge/patreon-donate-orange.svg)](https://www.patreon.com/octolab)
[![Build Status](https://travis-ci.org/kamilsk/retry.svg?branch=master)](https://travis-ci.org/kamilsk/retry)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamilsk/retry)](https://goreportcard.com/report/github.com/kamilsk/retry)
[![Coverage Status](https://coveralls.io/repos/github/kamilsk/retry/badge.svg)](https://coveralls.io/github/kamilsk/retry)
[![GoDoc](https://godoc.org/github.com/kamilsk/retry?status.svg)](https://godoc.org/github.com/kamilsk/retry)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](LICENSE)

## Differences from [Rican7/retry](https://github.com/Rican7/retry)

- Fixed [bug](https://github.com/Rican7/retry/pull/2) with an unexpected infinite loop.
  - Added a clear mechanism for this purpose as the Infinite [strategy](strategy/strategy.go#L24-L28).
- Added `context` support to cancellation.
- Added `error` transmission between attempts.
  - Added `classifier` to handle them (see [classifier](classifier) package).
- Added CLI tool `retry` which provides functionality for repeating terminal commands (see [cmd/retry](cmd/retry)).

## Usage

### Console tool for command execution with retries

This example shows how to repeat console command until successful.

```bash
$ retry --infinite -timeout 10m -backoff=lin:500ms -- /bin/sh -c 'echo "trying..."; exit $((1 + RANDOM % 10 > 5))'
```

[![asciicast](https://asciinema.org/a/150367.png)](https://asciinema.org/a/150367)

See more details [here](cmd/retry).

### Interrupt execution

```go
interrupter := retry.Multiplex(
	retry.WithTimeout(time.Second),
	retry.WithSignal(os.Interrupt),
)
if err := retry.Retry(interrupter, func(uint) error { time.Sleep(time.Second); return nil }); err == nil {
	panic("press Ctrl+C")
}
// successful interruption
```

## Installation

```bash
$ go get github.com/kamilsk/retry
```

### Mirror

```bash
$ egg bitbucket.org/kamilsk/retry
```

> [egg](https://github.com/kamilsk/egg) is an `extended go get`.

### Update

This library is using [SemVer](http://semver.org) for versioning, and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe. Therefore, do not use `go get -u` to update it,
use [dep](https://github.com/golang/dep) or something similar for this purpose.

## Feedback

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/kamilsk/retry)
[![@kamilsk](https://img.shields.io/badge/author-%40kamilsk-blue.svg)](https://twitter.com/ikamilsk)
[![@octolab](https://img.shields.io/badge/sponsor-%40octolab-blue.svg)](https://twitter.com/octolab_inc)

## Notes

- tested on Go 1.5, 1.6, 1.7, 1.8 and 1.9
- [research](../../tree/research)
- made with ❤️ by [OctoLab](https://www.octolab.org/)

[![Analytics](https://ga-beacon.appspot.com/UA-109817251-1/retry/dev)](https://github.com/igrigorik/ga-beacon)
