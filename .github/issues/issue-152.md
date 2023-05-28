---
code: null
id: MDU6SXNzdWU1OTA1MjUxMTc=
databaseId: 590525117
number: 152
url: https://github.com/kamilsk/retry/issues/152
title: "provide wrappers for retriable error"
issueType: null
labels:
  - "type:task"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-5]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2020-03-30T19:11:48Z
updatedAt: 2020-04-03T06:10:56Z
lastEditedAt: null
closedAt: 2020-04-03T06:10:56Z
---

# provide wrappers for retriable error

to simplify the code
https://github.com/goreleaser/goreleaser/commit/0736162d5efc262b324b9416b10750cd50182c41

```go
errors.Retriable{err}
```

<!-- 2020-03-30T20:02Z https://github.com/kamilsk/retry/issues/152#issuecomment-606215079
implement it here https://github.com/octolab/pkg/tree/master/errors
-->

<!-- 2020-04-03T06:10Z https://github.com/kamilsk/retry/issues/152#issuecomment-608249345
implemented
-->
