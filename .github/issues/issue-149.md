---
code: null
id: MDU6SXNzdWU1NzMzMzIxNDY=
databaseId: 573332146
number: 149
url: https://github.com/kamilsk/retry/issues/149
title: "refactor error check"
issueType: null
labels:
  - "type:refactoring"
assignees:
  - kamilsk
milestone: "[[milestone-6]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2020-02-29T16:45:53Z
updatedAt: 2020-03-01T18:00:48Z
lastEditedAt: null
closedAt: 2020-03-01T18:00:47Z
---

# refactor error check

```go
		strategy.CheckError(
			errors.Network(errors.Skip),
			errors.NoRows(errors.Skip),
		)
```
