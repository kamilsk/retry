---
code: null
id: MDU6SXNzdWUzODY0ODg3MjI=
databaseId: 386488722
number: 125
url: https://github.com/kamilsk/retry/issues/125
title: "remove strategy.Infinite"
issueType: null
labels:
  - "type:refactoring"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2018-12-01T19:26:06Z
updatedAt: 2019-02-17T14:22:42Z
lastEditedAt: null
closedAt: 2019-02-17T14:22:42Z
---

# remove strategy.Infinite

motivation:
- breaker as a first parameter
- goruntime.goshed on random iteration to prevent cpu bound
