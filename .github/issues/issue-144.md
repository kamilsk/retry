---
code: null
id: MDU6SXNzdWU1NjkzNzUxODQ=
databaseId: 569375184
number: 144
url: https://github.com/kamilsk/retry/issues/144
title: "restore error classifier"
issueType: null
labels:
  - "type:feature"
assignees:
  - kamilsk
milestone: "[[milestone-5]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2020-02-22T18:36:03Z
updatedAt: 2020-02-23T08:57:42Z
lastEditedAt: null
closedAt: 2020-02-23T08:57:42Z
---

# restore error classifier

```go
func CheckNetworkError(err error) bool {
	if err == nil {
		return true
	}
	if err, is := err.(net.Error); is {
		return err.Temporary()
	}
	return true
}
```
