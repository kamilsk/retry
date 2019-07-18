> # ♻️ retry
>
> The most advanced interruptible mechanism to perform actions repetitively until successful.

[![Build][icon_build]][page_build]
[![Quality][icon_quality]][page_quality]
[![Documentation][icon_docs]][page_docs]
[![Coverage][icon_coverage]][page_coverage]
[![Awesome][icon_awesome]][page_awesome]

## Usage

### Quick start

#### retry.Retry

```go
var response *http.Response

action := func(uint) error {
	var err error
	response, err = http.Get("https://github.com/kamilsk/retry")
	return err
}

if err := retry.Retry(breaker.BreakByTimeout(time.Minute), action, strategy.Limit(3)); err != nil {
	// handle error
}
// work with response
```

#### retry.Try

```go
var response *http.Response

action := func(uint) error {
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

if err := retry.Try(interrupter, action, strategy.Limit(3)); err != nil {
	// handle error
}
// work with response
```

or use Context

```go
ctx, cancel := context.WithTimeout(request.Context(), time.Minute)
defer cancel()

if err := retry.Try(ctx, action, strategy.Limit(3)); err != nil {
	// handle error
}
// work with response
```

#### retry.TryContext

```go
var response *http.Response

action := func(ctx context.Context, _ uint) error {
	req, err := http.NewRequest(http.MethodGet, "https://github.com/kamilsk/retry", nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	response, err = http.DefaultClient.Do(req)
	return err
}

// you can combine Context and Breaker together
interrupter, ctx := breaker.WithContext(request.Context())
defer interrupter.Close()

if err := retry.TryContext(ctx, action, strategy.Limit(3)); err != nil {
	// handle error
}
// work with response
```

#### Complex example

```go
what := func(uint) error {
	return do.Some("heavy work")
}
how := retry.How{
	strategy.Limit(3),
	strategy.BackoffWithJitter(
		backoff.Fibonacci(10*time.Millisecond),
		jitter.NormalDistribution(
			rand.New(rand.NewSource(time.Now().UnixNano())),
			0.25,
		),
	),
}
if err := retry.Try(ctx, what, how...); err != nil {
	log.Fatal(err)
}
```

## 🧩 Integration

The library uses [SemVer](https://semver.org) for versioning, and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe through major releases.
You can use [dep][] or [go modules][gomod] to manage its version.

The **[master][legacy]** is a feature frozen branch for versions **3.3.x** and no longer maintained.

```bash
$ dep ensure -add github.com/kamilsk/retry@3.3.3
```

The **[v3][]** branch is a continuation of the **[master][legacy]** branch for versions **v3.4.x**
to better integration with [go modules][gomod].

```bash
$ go get -u github.com/kamilsk/retry/v3@v3.4.4
```

The **[v4][]** branch is an actual development branch.

```bash
$ go get -u github.com/kamilsk/retry    # inside GOPATH and for old Go versions

$ go get -u github.com/kamilsk/retry/v4 # inside Go module, works well since Go 1.11

$ dep ensure -add github.com/kamilsk/retry@v4.0.0
```

Version **v4** focused on integration with the 🚧 [breaker][] package.

## 🤲 Outcomes

### Console tool for command execution with retries

This example shows how to repeat console command until successful.

```bash
$ retry -timeout 10m -backoff lin:500ms -- /bin/sh -c 'echo "trying..."; exit $((1 + RANDOM % 10 > 5))'
```

[![asciicast](https://asciinema.org/a/150367.png)](https://asciinema.org/a/150367)

See more details [here][cli].

---

made with ❤️ for everyone

[icon_awesome]:    https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg
[icon_build]:      https://travis-ci.org/kamilsk/retry.svg?branch=v4
[icon_coverage]:   https://api.codeclimate.com/v1/badges/ed88afbc0754e49e9d2d/test_coverage
[icon_docs]:       https://godoc.org/github.com/kamilsk/retry?status.svg
[icon_quality]:    https://goreportcard.com/badge/github.com/kamilsk/retry

[page_awesome]:    https://github.com/avelino/awesome-go#utilities
[page_build]:      https://travis-ci.org/kamilsk/retry
[page_coverage]:   https://codeclimate.com/github/kamilsk/retry/test_coverage
[page_docs]:       https://godoc.org/github.com/kamilsk/retry
[page_promo]:      https://github.com/kamilsk/retry
[page_quality]:    https://goreportcard.com/report/github.com/kamilsk/retry

[legacy]:          https://github.com/kamilsk/retry/tree/master
[v3]:              https://github.com/kamilsk/retry/tree/v3
[v4]:              https://github.com/kamilsk/retry/projects/4

[breaker]:         https://github.com/kamilsk/breaker
[cli]:             https://github.com/kamilsk/retry.cli
[dep]:             https://golang.github.io/dep/
[egg]:             https://github.com/kamilsk/egg
[gomod]:           https://github.com/golang/go/wiki/Modules
