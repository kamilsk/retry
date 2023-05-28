---
code: null
id: MDU6SXNzdWU1Njg1ODAzOTk=
databaseId: 568580399
number: 143
url: https://github.com/kamilsk/retry/issues/143
title: "create constant backoff"
issueType: null
labels:
  - "type:task"
assignees:
  - kamilsk
milestone: "[[milestone-5]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2020-02-20T21:22:16Z
updatedAt: 2020-02-21T19:52:35Z
lastEditedAt: null
closedAt: 2020-02-21T19:52:35Z
---

# create constant backoff

```go
func Constant(initial time.Duration) Algorithm {
  return func (...) time.Duration {
    return initial
  }
}
```
