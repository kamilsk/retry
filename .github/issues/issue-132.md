---
code: null
id: MDU6SXNzdWU0MjMxNDcwNTM=
databaseId: 423147053
number: 132
url: https://github.com/kamilsk/retry/issues/132
title: "go get v4 is not working"
issueType: null
labels:
  - "type:docs"
assignees:
  - kamilsk
milestone: null
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2019-03-20T09:25:34Z
updatedAt: 2019-09-10T12:03:21Z
lastEditedAt: null
closedAt: 2019-03-24T20:46:09Z
---

# go get v4 is not working

```
go get -u github.com/kamilsk/retry/v4
package github.com/kamilsk/retry/v4: cannot find package "github.com/kamilsk/retry/v4" in any of:
...
```

<!-- 2019-03-20T09:54Z https://github.com/kamilsk/retry/issues/132#issuecomment-474760660
Hi @mschneider82, it works only inside another package with enabled go modules.

```
$ mkdir my-new-package && cd $_
$ go mod init $_
$ go get -u github.com/kamilsk/retry/v4
go: finding github.com/kamilsk/retry/v4 v4.0.3
go: downloading github.com/kamilsk/retry/v4 v4.0.3
```

Thus it's not `go get`able anymore cause `go mod` is supported.

Maybe I have to add this into the docs. Thanks!
-->

<!-- 2019-03-20T10:14Z https://github.com/kamilsk/retry/issues/132#issuecomment-474767767
i tried to `go get -u github.com/goreleaser/goreleaser` and got this error, it works now
-->
