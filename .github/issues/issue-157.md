---
code: null
id: MDU6SXNzdWU2NzY1OTIxMzc=
databaseId: 676592137
number: 157
url: https://github.com/kamilsk/retry/issues/157
title: "strategy for context modification"
issueType: null
labels:
  - "type:feature"
assignees:
  - kamilsk
milestone: "[[milestone-6]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2020-08-11T05:28:56Z
updatedAt: 2021-01-08T12:39:51Z
lastEditedAt: null
closedAt: 2021-01-08T12:39:51Z
---

# strategy for context modification

```go
type Modifier func(context.Context) error
```

<!-- 2021-01-08T12:39Z https://github.com/kamilsk/retry/issues/157#issuecomment-756735310
won't do cause a3f6ed5
-->
