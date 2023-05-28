---
code: null
id: MDU6SXNzdWUyMjkwMDQ1Mjg=
databaseId: 229004528
number: 54
url: https://github.com/kamilsk/retry/issues/54
title: "callback with classifier"
issueType: null
labels:
  - "type:feature"
assignees:
  - kamilsk
milestone: null
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-05-16T11:50:57Z
updatedAt: 2017-05-24T17:30:55Z
lastEditedAt: null
closedAt: 2017-05-24T17:30:55Z
---

# callback with classifier

```go
err := retry.Retry(action, strategy.Callback(callback, classifier), ...)
```

`callback` is `func () err`, `classifier` handles `err` and makes a decision.
