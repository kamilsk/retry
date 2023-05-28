---
code: null
id: MDU6SXNzdWUyMjM2Mzk3MjI=
databaseId: 223639722
number: 50
url: https://github.com/kamilsk/retry/issues/50
title: "fix go 1.5 problems"
issueType: null
labels:
  - "type:critical"
assignees:
  - kamilsk
milestone: "[[milestone-1]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-04-23T14:31:30Z
updatedAt: 2017-04-26T15:48:40Z
lastEditedAt: null
closedAt: 2017-04-26T15:48:40Z
---

# fix go 1.5 problems

```
# github.com/kamilsk/retry/classifier
classifier/functional_test.go:43: unknown net.DNSError field 'IsTemporary' in struct literal
ok  	github.com/kamilsk/retry	1.113s
ok  	github.com/kamilsk/retry/backoff	1.010s
FAIL	github.com/kamilsk/retry/classifier [build failed]
?   	github.com/kamilsk/retry/cmd/retry	[no test files]
?   	github.com/kamilsk/retry/cmd/retry/flag	[no test files]
ok  	github.com/kamilsk/retry/jitter	1.015s
ok  	github.com/kamilsk/retry/strategy	1.171s
make: *** [docker-test-1.5] Error 123
```
