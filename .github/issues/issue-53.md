---
code: null
id: MDU6SXNzdWUyMjczNzYzODM=
databaseId: 227376383
number: 53
url: https://github.com/kamilsk/retry/issues/53
title: "retry cmd doesn't work properly"
issueType: null
labels:
  - "type:critical"
assignees:
  - kamilsk
milestone: null
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-05-09T13:58:28Z
updatedAt: 2017-05-25T10:40:50Z
lastEditedAt: 2017-05-25T10:40:48Z
closedAt: 2017-05-25T10:40:50Z
---

# retry cmd doesn't work properly

- [x] fix example

```bash
$ curl http://unknown.host
404 page not found
$ echo $?
0
```

- [x] fix execution

```bash
retry -backoff=lin[10s] --infinite -- curl
zsh: no matches found: -backoff=lin[10s]
```

```bash
$ retry --infinite -- exit 1
... nothing happen
```

- [x] fix goreleaser

```
-linkmode external -extldflags \"-static\"
```

https://github.com/goreleaser/goreleaser/pull/226/files#diff-39e929907e3737ff7460d982cb6fbc60R3

- [ ] do not panic

- [x] better logging
