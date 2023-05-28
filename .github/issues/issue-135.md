---
code: null
id: MDU6SXNzdWU0NzcyMTQyOTk=
databaseId: 477214299
number: 135
url: https://github.com/kamilsk/retry/issues/135
title: "TryContext doesn't work properly"
issueType: null
labels:
  - "type:bug"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2019-08-06T07:25:57Z
updatedAt: 2026-11-08T14:02:38Z
lastEditedAt: null
closedAt: 2019-08-07T20:39:01Z
---

# TryContext doesn't work properly

The main problem is that `ctx` is not interrupted by `func retry`.

Possible solution without breaking changes:

```go
func TryContext(
	ctx context.Context,
	action func(ctx context.Context, attempt uint) error,
	strategies ...func(attempt uint, err error) bool,
) error {
	ctx, cancel := context.WithCancel(ctx)
	err := retry(ctx, currying(ctx, action), strategies...)
	cancel() // or defer it
	return err
}
```

Also, I have to delete

```go
	if (breaker == nil || breaker.Done() == nil) && len(strategies) == 0 {
		return action(0)
	}
```

<!-- 2019-08-06T07:27Z https://github.com/kamilsk/retry/issues/135#issuecomment-518544103
Do not use defer by performance reason

```go
	if breaker != nil {
		defer breaker.Close()
	}
	return retry(breaker, action, strategies...)
```
-->
