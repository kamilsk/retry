---
code: null
id: MDU6SXNzdWUxODY5MTIyMDA=
databaseId: 186912200
number: 5
url: https://github.com/kamilsk/retry/issues/5
title: "fix gometalinter"
issueType: null
labels:
  - "type:refactoring"
assignees:
  - kamilsk
milestone: "[[milestone-1]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2016-11-02T20:32:52Z
updatedAt: 2026-11-08T14:00:45Z
lastEditedAt: 2026-11-08T14:00:45Z
closedAt: 2017-01-27T12:43:18Z
---

# fix gometalinter

```
$ make docker-gometalinter
docker run --rm \
	           -v '/Users/.../Development/go/src/github.com/kamilsk/retrier':'/go/src/github.com/kamilsk/retrier' \
	           -w '/go/src/github.com/kamilsk/retrier' \
	           kamilsk/go-tools:latest \
	           /bin/sh -c 'go list ./... | grep -v /vendor/ | xargs go test -i "$1" && \
	                       gometalinter.v1 --vendor --deadline=20s ./...'
retry.go:1:1:warning: package comment should be of the form "Package retrier ..." (golint)
jitter/jitter_test.go:9::warning: duplicate of jitter_test.go:29-47 (dupl)
jitter/jitter_test.go:29::warning: duplicate of jitter_test.go:9-27 (dupl)
strategy/extended_test.go:24:32:warning: error return value not checked (cannot use nil (untyped nil value) as parameter of type []Strategy) (errcheck)
strategy/extended_test.go:24:32:warning: unused global variable cannot use nil (untyped nil value) as parameter of type []Strategy (varcheck)
example_test.go:21:15:warning: error return value not checked (retrier.Retry(func(attempt uint) error {) (errcheck)
example_test.go:43:15:warning: error return value not checked (logFile.Chdir() // Do something with the file) (errcheck)
example_test.go:83:15:warning: error return value not checked (retrier.Retry() (errcheck)
example_test.go:98:11:warning: error return value not checked (rw.Write([]byte("Internal Server Error"))) (errcheck)
example_test.go:143:11:warning: error return value not checked (rw.Write([]byte("Internal Server Error"))) (errcheck)
make: *** [docker-tool-gometalinter] Error 1
```

<!-- 2017-01-27T10:16Z https://github.com/kamilsk/retry/issues/5#issuecomment-275632419
```
$ make docker-gometalinter
docker run --rm \
	           -v '/Users/.../Development/go/src/github.com/kamilsk/retrier':'/go/src/github.com/kamilsk/retrier' \
	           -w '/go/src/github.com/kamilsk/retrier' \
	           kamilsk/go-tools:latest \
	           /bin/sh -c 'go list ./... | grep -v /vendor/ | xargs go test -i && \
	                       gometalinter.v1 --vendor --deadline=1m ./...'
example_test.go:24:15:warning: error return value not checked (retrier.Retry(context.Background(), func(attempt uint) error {) (errcheck)
example_test.go:46:15:warning: error return value not checked (logFile.Chdir() // Do something with the file) (errcheck)
example_test.go:87:15:warning: error return value not checked (retrier.Retry() (errcheck)
example_test.go:124:11:warning: error return value not checked (rw.Write([]byte("Internal Server Error"))) (errcheck)
make: *** [docker-tool-gometalinter] Error 1
```
-->
