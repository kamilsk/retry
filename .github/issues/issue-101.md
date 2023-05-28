---
code: null
id: MDU6SXNzdWUyODIzMzM4MDI=
databaseId: 282333802
number: 101
url: https://github.com/kamilsk/retry/issues/101
title: "WithContext should return ctx.Done channel"
issueType: null
labels:
  - "type:task"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-12-15T06:39:51Z
updatedAt: 2019-02-17T09:40:35Z
lastEditedAt: null
closedAt: 2019-02-17T09:40:35Z
---

# WithContext should return ctx.Done channel

```go
semaphore.Acquire(semaphore.Multiplex(semaphore.WithContext(req.Context()), semaphore.WithSignal(os.Interrupt)))
```

like as errgroup:

```go
g, ctx := errgroup.WithContext(ctx)
```

https://github.com/golang/sync/tree/master/errgroup

<!-- 2019-02-17T09:40Z https://github.com/kamilsk/retry/issues/101#issuecomment-464434764
won't do
-->
