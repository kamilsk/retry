---
code: null
id: MDU6SXNzdWU1NjU1MTA1ODI=
databaseId: 565510582
number: 142
url: https://github.com/kamilsk/retry/issues/142
title: "review Strategy type"
issueType: null
labels:
  - "type:feature"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-6]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2020-02-14T19:18:17Z
updatedAt: 2020-02-16T19:11:47Z
lastEditedAt: null
closedAt: 2020-02-16T19:11:47Z
---

# review Strategy type

```go
type Strategy func(Iteration) bool
```

```go
type Iteration struct {
    attempt uint
    breaker Breaker
    error error
}
```
