> # 🚧 breaker
>
> Flexible mechanism to make your code breakable.

[![Patreon][icon_patreon]][support]
[![GoDoc][icon_docs]][docs]
[![License][icon_license]][license]

A Breaker carries a cancellation signal to break an action execution.

Example based on [github.com/kamilsk/retry][retry] package:

```go
if err := retry.Retry(breaker.BreakByTimeout(time.Minute), action); err != nil {
	log.Fatal(err)
}
```

Example based on [github.com/kamilsk/semaphore][semaphore] package:

```go
if err := semaphore.Acquire(breaker.BreakByTimeout(time.Minute), 5); err != nil {
	log.Fatal(err)
}
```

Complex example:

```go
interrupter := breaker.Multiplex(
	func () breaker.Interface {
		br, _ := breaker.WithContext(request.Context())
		return br
	}()
	breaker.BreakByTimeout(time.Minute),
	breaker.BreakBySignal(os.Interrupt),
)
defer interrupter.Close()

<-interrupter.Done() // wait context cancellation, timeout or interrupt signal
```

## Notice

This package is based on the [platform][] - my toolset to build microservices such as [click][] or [forma][].
It is stable, well-tested and production ready.

---

[![@kamilsk][icon_tw_author]][author]
[![@octolab][icon_tw_sponsor]][sponsor]

made with ❤️ by [OctoLab][octolab]

[docs]:            https://godoc.org/github.com/kamilsk/breaker
[license]:         LICENSE
[promo]:           https://github.com/kamilsk/breaker

[click]:           https://github.com/kamilsk/click
[forma]:           https://github.com/kamilsk/form-api
[platform]:        https://github.com/kamilsk/platform
[retry]:           https://github.com/kamilsk/retry
[semaphore]:       https://github.com/kamilsk/semaphore

[author]:          https://twitter.com/ikamilsk
[octolab]:         https://www.octolab.org/
[sponsor]:         https://twitter.com/octolab_inc
[support]:         https://www.patreon.com/octolab

[icon_docs]:       https://godoc.org/github.com/kamilsk/breaker?status.svg
[icon_license]:    https://img.shields.io/badge/license-MIT-blue.svg
[icon_patreon]:    https://img.shields.io/badge/patreon-donate-orange.svg
[icon_tw_author]:  https://img.shields.io/badge/author-%40kamilsk-blue.svg
[icon_tw_sponsor]: https://img.shields.io/badge/sponsor-%40octolab-blue.svg
