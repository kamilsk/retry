---
code: null
id: MDU6SXNzdWU3Nzc2MjAyNDM=
databaseId: 777620243
number: 161
url: https://github.com/kamilsk/retry/issues/161
title: "fix lint issues and coverage report"
issueType: null
labels:
  - "type:task"
assignees:
  - kamilsk
milestone: "[[milestone-6]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2021-01-03T11:33:41Z
updatedAt: 2021-01-08T15:31:33Z
lastEditedAt: null
closedAt: 2021-01-08T15:31:32Z
---

# fix lint issues and coverage report

```
example_test.go:42:3: exitAfterDefer: log.Fatal clutters `defer cancel()` (gocritic)
		log.Fatal(err)
		^
retry_test.go:65:44: SA1029: should not use built-in type string as key for value; define your own type to avoid collisions (staticcheck)
		ctx := context.WithValue(context.TODO(), "test", "value")
		                                         ^
retry_test.go:140:44: SA1029: should not use built-in type string as key for value; define your own type to avoid collisions (staticcheck)
		ctx := context.WithValue(context.TODO(), "test", "value")
		                                         ^
```

+ https://travis-ci.com/github/kamilsk/retry/jobs/467785097

```
$ if [[ $TRAVIS_GO_VERSION == 1.15* ]]; then
    cc-test-reporter after-build -t gocov -p $(go list -m) --exit-code $TRAVIS_TEST_RESULT
fi

time="2021-01-03T11:26:18Z" level=error msg="failed to read file v5/backoff/backoff.go\nopen v5/backoff/backoff.go: no such file or directory"
Error: open v5/backoff/backoff.go: no such file or directory
```
