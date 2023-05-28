---
code: null
id: MDU6SXNzdWUyMjA0NzUxMzI=
databaseId: 220475132
number: 41
url: https://github.com/kamilsk/retry/issues/41
title: "prevent output to stderr when context has closed"
issueType: null
labels:
  - "type:bug"
assignees:
  - kamilsk
milestone: "[[milestone-1]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-04-09T14:37:02Z
updatedAt: 2017-04-15T08:58:16Z
lastEditedAt: null
closedAt: 2017-04-15T08:58:16Z
---

# prevent output to stderr when context has closed

```
$ make cmd-test
go install -ldflags "-X 'main.Timeout=100ms'" ./cmd/retry
retry -limit=3 -backoff=lin[10ms] -timeout=500ms -- curl http://unknown.host
error occurred "context deadline exceeded"
curl: (7) Failed to connect to unknown.host port 80: Connection refused
```
