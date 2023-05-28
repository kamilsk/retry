---
code: null
id: MDU6SXNzdWU1Njk2OTk1NDA=
databaseId: 569699540
number: 146
url: https://github.com/kamilsk/retry/issues/146
title: "extend error check strategies"
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
createdAt: 2020-02-24T08:20:22Z
updatedAt: 2020-02-24T09:17:18Z
lastEditedAt: null
closedAt: 2020-02-24T09:17:18Z
---

# extend error check strategies

```go
CheckError

Retriable() bool

const (
  Skip = true
  Strict = false
)

CheckNetworkError(default bool)
```
