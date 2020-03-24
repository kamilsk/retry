> # ♻️ retry
>
> The most advanced interruptible mechanism to perform actions repetitively until successful.

[![Build][build.icon]][build.page]
[![Documentation][docs.icon]][docs.page]
[![Quality][quality.icon]][quality.page]
[![Template][template.icon]][template.page]
[![Coverage][coverage.icon]][coverage.page]
[![Awesome][awesome.icon]][awesome.page]

## 💡 Idea

The package based on [Rican7/retry][] but fully reworked and focused on integration
with the 🚧 [breaker][] and the built-in [context][] packages.

Full description of the idea is available [here][design.page].

## 🏆 Motivation

I developed distributed systems at [Lazada][], and later at [Avito][],
which communicate with each other through a network, and I need a package to make
these communications more reliable.

## 🤼‍♂️ How to

### retry.Do

```go
var response *http.Response

action := func() error {
	var err error
	response, err = http.Get("https://github.com/kamilsk/retry")
	return err
}

// you can combine multiple Breakers into one
interrupter := breaker.MultiplexTwo(
	breaker.BreakByTimeout(time.Minute),
	breaker.BreakBySignal(os.Interrupt),
)
defer interrupter.Close()

if err := retry.Do(interrupter, action, strategy.Limit(3)); err != nil {
	if err == breaker.Interrupted {
		// operation was interrupted
	}
	// handle error
}
// work with response
```

or use Context

```go
ctx, cancel := context.WithTimeout(request.Context(), time.Minute)
defer cancel()

if err := retry.Do(ctx, action, strategy.Limit(3)); err != nil {
	if err == context.Canceled || err == context.DeadlineExceeded {
		// operation was interrupted
	}
	// handle error
}
// work with response
```

### Complex example

```go
import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/kamilsk/retry/v5"
	"github.com/kamilsk/retry/v5/backoff"
	"github.com/kamilsk/retry/v5/jitter"
	"github.com/kamilsk/retry/v5/strategy"
)

func main() {
	what := SendRequest

	how := retry.How{
		strategy.Limit(5),
		strategy.BackoffWithJitter(
			backoff.Fibonacci(10*time.Millisecond),
			jitter.NormalDistribution(
				rand.New(rand.NewSource(time.Now().UnixNano())),
				0.25,
			),
		),
		strategy.CheckError(
			strategy.NetworkError(strategy.Skip),
			DatabaseError(),
		),
	}

	breaker, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := retry.Do(breaker, what, how...); err != nil {
		log.Fatal(err)
	}
}

func SendRequest() error {
	// communicate with some service
}

func DatabaseError() func(error) bool {
	blacklist := []error{sql.ErrNoRows, sql.ErrConnDone, sql.ErrTxDone}
	return func(err error) bool {
		for _, preset := range blacklist {
			if err == preset {
				return false
			}
		}
		return true
	}
}
```

## 🧩 Integration

The library uses [SemVer](https://semver.org) for versioning, and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe through major releases.
You can use [go modules](https://github.com/golang/go/wiki/Modules) to manage its version.

```bash
$ go get github.com/kamilsk/retry/v5@latest
```

## 🤲 Outcomes

### Console tool for command execution with retries

This example shows how to repeat console command until successful.

```bash
$ retry -timeout 10m -backoff lin:500ms -- /bin/sh -c 'echo "trying..."; exit $((1 + RANDOM % 10 > 5))'
```

[![asciicast][cli.preview]][cli.demo]

See more details [here][cli].

---

made with ❤️ for everyone

[awesome.page]:     https://github.com/avelino/awesome-go#utilities
[awesome.icon]:     https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg
[build.page]:       https://travis-ci.org/kamilsk/retry
[build.icon]:       https://travis-ci.org/kamilsk/retry.svg?branch=v5
[coverage.page]:    https://codeclimate.com/github/kamilsk/retry/test_coverage
[coverage.icon]:    https://api.codeclimate.com/v1/badges/ed88afbc0754e49e9d2d/test_coverage
[design.page]:      https://www.notion.so/octolab/retry-cab5722faae445d197e44fbe0225cc98?r=0b753cbf767346f5a6fd51194829a2f3
[docs.page]:        https://pkg.go.dev/github.com/kamilsk/retry/v5
[docs.icon]:        https://img.shields.io/badge/docs-pkg.go.dev-blue
[promo.page]:       https://github.com/kamilsk/retry
[quality.page]:     https://goreportcard.com/report/github.com/kamilsk/retry
[quality.icon]:     https://goreportcard.com/badge/github.com/kamilsk/retry
[template.page]:    https://github.com/octomation/go-module
[template.icon]:    https://img.shields.io/badge/template-go--module-blue

[Avito]:            https://tech.avito.ru
[breaker]:          https://github.com/kamilsk/breaker
[cli]:              https://github.com/kamilsk/try
[cli.demo]:         https://asciinema.org/a/150367
[cli.preview]:      https://asciinema.org/a/150367.png
[context]:          https://pkg.go.dev/context
[Lazada]:           https://github.com/lazada
[Rican7/retry]:     https://github.com/Rican7/retry

[tmp.docs]:         https://nicedoc.io/kamilsk/retry?theme=dark
[tmp.history]:      https://github.githistory.xyz/kamilsk/retry/blob/v5/README.md
