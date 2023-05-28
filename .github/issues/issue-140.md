---
code: null
id: MDU6SXNzdWU1NjU1MDM5NTQ=
databaseId: 565503954
number: 140
url: https://github.com/kamilsk/retry/issues/140
title: "extend breaker interface"
issueType: null
labels:
  - "type:feature"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-6]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2020-02-14T19:06:33Z
updatedAt: 2020-02-16T08:26:37Z
lastEditedAt: null
closedAt: 2020-02-16T08:26:36Z
---

# extend breaker interface

```
// A Breaker carries a cancellation signal to break an action execution.
//
// It is a subset of context.Context and github.com/kamilsk/breaker.Breaker.
type Breaker interface {
	// Done returns a channel that's closed when a cancellation signal occurred.
	Done() <-chan struct{}

	Error() error
}
```

<!-- 2020-02-15T09:14Z https://github.com/kamilsk/retry/issues/140#issuecomment-586570944
blocked by https://github.com/kamilsk/breaker/issues/8
-->
