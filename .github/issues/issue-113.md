---
code: null
id: MDU6SXNzdWUyOTgxMDcwMjQ=
databaseId: 298107024
number: 113
url: https://github.com/kamilsk/retry/issues/113
title: "race in examples"
issueType: null
labels:
  - "type:bug"
assignees:
  - kamilsk
milestone: "[[milestone-3]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2018-02-18T18:38:46Z
updatedAt: 2026-11-08T14:00:37Z
lastEditedAt: 2026-11-08T14:00:37Z
closedAt: 2018-02-18T19:20:15Z
---

# race in examples

```
make test
go list ./... | grep -v vendor | grep -v ^_ | xargs go test -race
==================
WARNING: DATA RACE
Read at 0x00c42008e4d8 by goroutine 62:
  github.com/kamilsk/retry_test.Example_httpClient.func1()
      /Users/.../Development/go/src/github.com/kamilsk/retry/example_http_client_test.go:43 +0x4a
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.10/libexec/src/net/http/server.go:1947 +0x51
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.10/libexec/src/net/http/server.go:2694 +0xb9
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.10/libexec/src/net/http/server.go:1830 +0x7dc

Previous write at 0x00c42008e4d8 by goroutine 57:
  github.com/kamilsk/retry_test.Example_httpClient.func1()
      /Users/.../Development/go/src/github.com/kamilsk/retry/example_http_client_test.go:47 +0xd9
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.10/libexec/src/net/http/server.go:1947 +0x51
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.10/libexec/src/net/http/server.go:2694 +0xb9
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.10/libexec/src/net/http/server.go:1830 +0x7dc

Goroutine 62 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.10/libexec/src/net/http/server.go:2795 +0x364
  net/http/httptest.(*Server).goServe.func1()
      /usr/local/Cellar/go/1.10/libexec/src/net/http/httptest/server.go:280 +0xa2

Goroutine 57 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.10/libexec/src/net/http/server.go:2795 +0x364
  net/http/httptest.(*Server).goServe.func1()
      /usr/local/Cellar/go/1.10/libexec/src/net/http/httptest/server.go:280 +0xa2
==================
FAIL
FAIL    github.com/kamilsk/retry        0.685s
ok      github.com/kamilsk/retry/backoff        1.024s
ok      github.com/kamilsk/retry/classifier     1.029s
ok      github.com/kamilsk/retry/cmd/retry      1.047s
ok      github.com/kamilsk/retry/jitter 1.024s
ok      github.com/kamilsk/retry/strategy       1.185s
```
