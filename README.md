> # ♻️ retry
>
> Functional mechanism based on channels to perform actions repetitively until successful.

[![Awesome][icon_awesome]][awesome]
[![Patreon][icon_patreon]][support]
[![Build Status][icon_build]][build]
[![Code Coverage][icon_coverage]][quality]
[![Code Quality][icon_quality]][quality]
[![GoDoc][icon_docs]][docs]
[![Research][icon_research]][research]
[![License][icon_license]][license]

## Important news

The **[master][legacy]** is a feature frozen branch for versions _3.3.x_ and no longer maintained.

```bash
$ dep ensure -add github.com/kamilsk/retry@3.3.2
```

The **[v3][]** branch is a continuation of the **master** branch for versions _v3.4.y_
to better integration with [Go Modules][gomod].

```bash
$ go get -u github.com/kamilsk/retry/v3@v3.4.3
```

The **v4** branch is an actual development branch with many [features][v4_features].

```bash
$ go get -u github.com/kamilsk/retry/v4
```

## Usage

### Quick start

```go
var (
	response *http.Response
	action   retry.Action = func(_ uint) error {
		var err error
		response, err = http.Get("https://github.com/kamilsk/retry")
		return err
	}
)

if err := retry.Retry(retry.WithTimeout(time.Minute), action, strategy.Limit(10)); err != nil {
	// handle error
	return
}
// work with response
```

### Console tool for command execution with retries

This example shows how to repeat console command until successful.

```bash
$ retry --infinite -timeout 10m -backoff=lin:500ms -- /bin/sh -c 'echo "trying..."; exit $((1 + RANDOM % 10 > 5))'
```

[![asciicast](https://asciinema.org/a/150367.png)](https://asciinema.org/a/150367)

See more details [here](cmd/retry).

### Create HTTP client with retry

This example shows how to extend standard http.Client with retry under the hood.

```go
type client struct {
	base       *http.Client
	strategies []strategy.Strategy
}

func New(timeout time.Duration, strategies ...strategy.Strategy) *client {
	return &client{
		base:       &http.Client{Timeout: timeout},
		strategies: strategies,
	}
}

func (c *client) Get(deadline <-chan struct{}, url string) (*http.Response, error) {
	var response *http.Response
	err := retry.Retry(deadline, func(uint) error {
		resp, err := c.base.Get(url)
		if err != nil {
			return err
		}
		response = resp
		return nil
	}, c.strategies...)
	return response, err
}
```

### Control database connection

This example shows how to use retry to restore database connection by `database/sql/driver.Pinger`.

```go
MustOpen := func() *sql.DB {
	db, err := sql.Open("stub", "stub://test")
	if err != nil {
		panic(err)
	}
	return db
}

go func(db *sql.DB, ctx context.Context, shutdown chan<- struct{}, frequency time.Duration,
	strategies ...strategy.Strategy) {

	defer func() {
		if r := recover(); r != nil {
			shutdown <- struct{}{}
		}
	}()

	ping := func(uint) error {
		return db.Ping()
	}

	for {
		if err := retry.Retry(ctx.Done(), ping, strategies...); err != nil {
			panic(err)
		}
		time.Sleep(frequency)
	}
}(MustOpen(), context.Background(), shutdown, time.Millisecond, strategy.Limit(1))
```

### Use context for cancellation

This example shows how to use context and retry together.

```go
communication := make(chan error)

go service.Listen(communication)

action := func(uint) error {
	communication <- nil   // ping
	return <-communication // pong
}
ctx := retry.WithContext(context.Background(), retry.WithTimeout(time.Second))
if err := retry.Retry(ctx.Done(), action, strategy.Delay(time.Millisecond)); err != nil {
	// the service does not respond within one second
}
```

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
$ # or use mirror
$ egg bitbucket.org/kamilsk/retry
```

> [egg][]<sup id="anchor-egg">[1](#egg)</sup> is an `extended go get`.

## Update

This library is using [SemVer](https://semver.org/) for versioning, and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe. Therefore, do not use `go get -u` to update it,
use **dep**, **glide** or something similar for this purpose.

<sup id="egg">1</sup> The project is still in prototyping. [↩](#anchor-egg)

---

[![Gitter][icon_gitter]][gitter]
[![@kamilsk][icon_tw_author]][author]
[![@octolab][icon_tw_sponsor]][sponsor]

made with ❤️ by [OctoLab][octolab]

[awesome]:         https://github.com/avelino/awesome-go#utilities
[build]:           https://travis-ci.org/kamilsk/retry
[docs]:            https://godoc.org/github.com/kamilsk/retry
[gitter]:          https://gitter.im/kamilsk/retry
[license]:         LICENSE
[promo]:           https://github.com/kamilsk/retry
[quality]:         https://scrutinizer-ci.com/g/kamilsk/retry/?branch=v4
[research]:        https://github.com/kamilsk/go-research/tree/master/projects/retry
[legacy]:          https://github.com/kamilsk/retry/tree/master
[v3]:              https://github.com/kamilsk/retry/tree/v3
[v4_features]:     https://github.com/kamilsk/retry/projects/4

[egg]:             https://github.com/kamilsk/egg
[gomod]:           https://github.com/golang/go/wiki/Modules

[author]:          https://twitter.com/ikamilsk
[octolab]:         https://www.octolab.org/
[sponsor]:         https://twitter.com/octolab_inc
[support]:         https://www.patreon.com/octolab

[analytics]:       https://ga-beacon.appspot.com/UA-109817251-1/retry/v4?pixel
[tweet]:           https://twitter.com/intent/tweet?text=Functional%20mechanism%20based%20on%20channels%20to%20perform%20actions%20repetitively%20until%20successful&url=https://github.com/kamilsk/retry&via=ikamilsk&hashtags=go,repeat,retry,backoff,jitter

[icon_awesome]:    https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg
[icon_build]:      https://travis-ci.org/kamilsk/retry.svg?branch=v4
[icon_coverage]:   https://scrutinizer-ci.com/g/kamilsk/retry/badges/coverage.png?b=v4
[icon_docs]:       https://godoc.org/github.com/kamilsk/retry?status.svg
[icon_gitter]:     https://badges.gitter.im/Join%20Chat.svg
[icon_license]:    https://img.shields.io/badge/license-MIT-blue.svg
[icon_patreon]:    https://img.shields.io/badge/patreon-donate-orange.svg
[icon_quality]:    https://scrutinizer-ci.com/g/kamilsk/retry/badges/quality-score.png?b=v4
[icon_research]:   https://img.shields.io/badge/research-in%20progress-yellow.svg
[icon_tw_author]:  https://img.shields.io/badge/author-%40kamilsk-blue.svg
[icon_tw_sponsor]: https://img.shields.io/badge/sponsor-%40octolab-blue.svg
[icon_twitter]:    https://img.shields.io/twitter/url/http/shields.io.svg?style=social
