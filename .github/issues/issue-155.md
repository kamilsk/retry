---
code: null
id: MDU6SXNzdWU2MDU1NTgyMzc=
databaseId: 605558237
number: 155
url: https://github.com/kamilsk/retry/issues/155
title: "goroutine leak in v4"
issueType: null
labels:
  - "type:critical"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-6]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2020-04-23T13:54:22Z
updatedAt: 2020-04-30T19:22:53Z
lastEditedAt: null
closedAt: 2020-04-30T19:22:52Z
---

# goroutine leak in v4

In v4 it is possible to use "infinite" context to call inside an action

```
	err = retry.Try(
		breaker.BreakByTimeout(10*time.Second),
		func(uint) error {
			statuses, err = s.dwhClient.V1Status(ctx, source)
			return err
		},
		strategy.Limit(10),
		strategy.Wait(time.Millisecond*10, time.Millisecond*100, time.Millisecond*1000),
	)
	return
```

The problem is that if the client doesn't have read/conn timeout the V1Status can block goroutine for a long time.

<!-- 2020-04-23T13:57Z https://github.com/kamilsk/retry/issues/155#issuecomment-618409128
The possible solution:

```go
type Action func(context.Context) error
```
-->
