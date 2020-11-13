> # ♻️ retry [![Awesome][awesome.icon]][awesome.page]
>
> The most advanced interruptible mechanism to perform actions repetitively until successful.

[![Build][build.icon]][build.page]
[![Documentation][docs.icon]][docs.page]
[![Quality][quality.icon]][quality.page]
[![Template][template.icon]][template.page]
[![Coverage][coverage.icon]][coverage.page]
[![Mirror][mirror.icon]][mirror.page]

## 💡 Idea

The package based on [Rican7/retry][] but fully reworked and focused on integration
with the 🚧 [breaker][] and the built-in [context][] packages.

A full description of the idea is available [here][design.page].

## 🏆 Motivation

I developed distributed systems at [Lazada][], and later at [Avito][],
which communicate with each other through a network, and I need a package to make
these communications more reliable.

## 🤼‍♂️ How to

rewriting...

## 🧩 Integration

The library uses [SemVer](https://semver.org) for versioning, and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe through major releases.
You can use [go modules](https://github.com/golang/go/wiki/Modules) to manage its version.

```bash
$ go get github.com/kamilsk/retry/v5@latest
```

---

made with ❤️ for everyone

[build.page]:       https://travis-ci.com/kamilsk/retry
[build.icon]:       https://travis-ci.com/kamilsk/retry.svg?branch=v5
[coverage.page]:    https://codeclimate.com/github/kamilsk/retry/test_coverage
[coverage.icon]:    https://api.codeclimate.com/v1/badges/ed88afbc0754e49e9d2d/test_coverage
[design.page]:      https://www.notion.so/octolab/retry-cab5722faae445d197e44fbe0225cc98?r=0b753cbf767346f5a6fd51194829a2f3
[docs.page]:        https://pkg.go.dev/github.com/kamilsk/retry/v5
[docs.icon]:        https://img.shields.io/badge/docs-pkg.go.dev-blue
[promo.page]:       https://github.com/kamilsk/retry
[quality.page]:     https://goreportcard.com/report/github.com/kamilsk/retry/v5
[quality.icon]:     https://goreportcard.com/badge/github.com/kamilsk/retry/v5
[template.page]:    https://github.com/octomation/go-module
[template.icon]:    https://img.shields.io/badge/template-go--module-blue
[mirror.page]:      https://bitbucket.org/kamilsk/retry
[mirror.icon]:      https://img.shields.io/badge/mirror-bitbucket-blue

[awesome.page]:     https://github.com/avelino/awesome-go#utilities
[awesome.icon]:     https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg

[Avito]:            https://tech.avito.ru
[breaker]:          https://github.com/kamilsk/breaker
[cli]:              https://github.com/octolab/try
[context]:          https://pkg.go.dev/context
[Lazada]:           https://github.com/lazada
[Rican7/retry]:     https://github.com/Rican7/retry
