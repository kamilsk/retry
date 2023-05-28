---
code: null
id: MDU6SXNzdWUyODIzMzU5MzI=
databaseId: 282335932
number: 102
url: https://github.com/kamilsk/retry/issues/102
title: "add streaming"
issueType: null
labels:
  - "type:feature"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-12-15T06:53:43Z
updatedAt: 2019-02-23T20:08:07Z
lastEditedAt: 2017-12-15T06:53:52Z
closedAt: 2019-02-23T20:08:07Z
---

# add streaming

```go
retry --stream-to=12345 -- do command
```

retry starts listen 12345 TCP port to stream stdout/stderr of executed command to it.

<!-- 2017-12-15T06:54Z https://github.com/kamilsk/retry/issues/102#issuecomment-351928646
can be useful https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
-->

<!-- 2018-04-01T18:00Z https://github.com/kamilsk/retry/issues/102#issuecomment-377804883
some inspiration:
- https://github.com/wagoodman/bashful
-->

<!-- 2019-02-23T20:08Z https://github.com/kamilsk/retry/issues/102#issuecomment-466687864
moved to https://github.com/kamilsk/retry.cli/issues/4
-->
