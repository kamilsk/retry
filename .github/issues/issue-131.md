---
code: null
id: MDU6SXNzdWU0MTMzMTY0ODI=
databaseId: 413316482
number: 131
url: https://github.com/kamilsk/retry/issues/131
title: "find way to support dep for v4"
issueType: null
labels:
  - "type:critical"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2019-02-22T09:32:04Z
updatedAt: 2019-02-22T11:50:00Z
lastEditedAt: null
closedAt: 2019-02-22T11:50:00Z
---

# find way to support dep for v4

```
$ dep ensure -add github.com/kamilsk/retry/v4@v4.0.0
failed

$ dep ensure -add github.com/kamilsk/retry@v4.0.0
failed
```
