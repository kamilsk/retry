---
code: null
id: MDU6SXNzdWUyODMxNDYwMjM=
databaseId: 283146023
number: 106
url: https://github.com/kamilsk/retry/issues/106
title: "fix code quality problems"
issueType: null
labels:
  - "type:critical"
  - "type:refactoring"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-12-19T08:48:45Z
updatedAt: 2019-02-17T14:25:11Z
lastEditedAt: null
closedAt: 2019-02-17T14:25:11Z
---

# fix code quality problems

```
cmd/retry/main.go:116:14:warning: error return value not checked (spin.Color("red")) (errcheck)
cmd/retry/main.go:117:37:warning: error return value not checked (color.New(color.FgYellow).Fprintf(stderr, "#%d attempt at %s... \n", attempt+1, time.Now().Sub(start))) (errcheck)
cmd/retry/main.go:119::warning: Subprocess launching with variable.,HIGH,HIGH (gas)
cmd/retry/main.go:67:34:warning: error return value not checked (color.New(color.FgRed).Fprintf(app.Stderr, format, err)) (errcheck)
cmd/retry/main.go:88:39:warning: error return value not checked (color.New(color.FgYellow).Fprintln(stderr, "notify component is not ready yet")) (errcheck)
cmd/retry/main.go:90:18:warning: error return value not checked (report.Execute(app.Stdout, struct {) (errcheck)
cmd/retry/parser_test.go:13::warning: cyclomatic complexity 11 of function Test_parse() is high (> 10) (gocyclo)
```

<!-- 2017-12-19T08:52Z https://github.com/kamilsk/retry/issues/106#issuecomment-352678624
also fix

```
$ make test
go list ./... | grep -v vendor | grep -v ^_ | xargs go test -race
==================
WARNING: DATA RACE
Read at 0x00c420088540 by goroutine 59:
  github.com/kamilsk/retry_test.Example_httpClient.func1()
      /Users/.../Development/go/src/github.com/kamilsk/retry/example_http_client_test.go:43 +0x4d
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.9.2/libexec/src/net/http/server.go:1918 +0x51
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.9.2/libexec/src/net/http/server.go:2619 +0xbc
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.9.2/libexec/src/net/http/server.go:1801 +0x83b

Previous write at 0x00c420088540 by goroutine 54:
  github.com/kamilsk/retry_test.Example_httpClient.func1()
      /Users/.../Development/go/src/github.com/kamilsk/retry/example_http_client_test.go:47 +0xe2
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.9.2/libexec/src/net/http/server.go:1918 +0x51
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.9.2/libexec/src/net/http/server.go:2619 +0xbc
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.9.2/libexec/src/net/http/server.go:1801 +0x83b

Goroutine 59 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.9.2/libexec/src/net/http/server.go:2720 +0x37c
  net/http/httptest.(*Server).goServe.func1()
      /usr/local/Cellar/go/1.9.2/libexec/src/net/http/httptest/server.go:280 +0xa2

Goroutine 54 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.9.2/libexec/src/net/http/server.go:2720 +0x37c
  net/http/httptest.(*Server).goServe.func1()
      /usr/local/Cellar/go/1.9.2/libexec/src/net/http/httptest/server.go:280 +0xa2
==================
FAIL
FAIL    github.com/kamilsk/retry        0.678s
ok      github.com/kamilsk/retry/backoff        1.030s
ok      github.com/kamilsk/retry/classifier     1.024s
ok      github.com/kamilsk/retry/cmd/retry      1.050s
ok      github.com/kamilsk/retry/jitter 1.027s
ok      github.com/kamilsk/retry/strategy       1.191s
make: *** [test] Error 1
```
-->

<!-- 2019-02-17T14:25Z https://github.com/kamilsk/retry/issues/106#issuecomment-464463397
fixed
-->
