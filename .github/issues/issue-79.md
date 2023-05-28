---
code: null
id: MDU6SXNzdWUyNDY1OTQ3NDE=
databaseId: 246594741
number: 79
url: https://github.com/kamilsk/retry/issues/79
title: "add example with os.Interrupt"
issueType: null
labels:
  - "type:docs"
assignees:
  - kamilsk
milestone: "[[milestone-3]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-07-30T16:43:15Z
updatedAt: 2017-11-30T14:47:53Z
lastEditedAt: null
closedAt: 2017-11-30T14:47:53Z
---

# add example with os.Interrupt

concept
```go
	ctx := context.Background()

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()
```
