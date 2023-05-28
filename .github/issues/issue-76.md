---
code: null
id: MDU6SXNzdWUyNDQ4OTAxNTA=
databaseId: 244890150
number: 76
url: https://github.com/kamilsk/retry/issues/76
title: "add strategy based on classifier package"
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
createdAt: 2017-07-23T05:33:44Z
updatedAt: 2019-02-17T07:54:56Z
lastEditedAt: 2017-11-30T15:13:52Z
closedAt: 2019-02-17T07:54:56Z
---

# add strategy based on classifier package

the concept:
```go
func Whitelist(errors ...error) Strategy {
	whitelist := classifier.WhitelistClassifier(errors)
	return func(attempt uint, err error) bool {
		return whitelist.Classify(err) == classifier.Retry
	}
}
```

<!-- 2017-12-14T09:37Z https://github.com/kamilsk/retry/issues/76#issuecomment-351658240
also, add the example from `example_error_test.go`
-->

<!-- 2019-02-17T07:54Z https://github.com/kamilsk/retry/issues/76#issuecomment-464426529
won't do
-->
