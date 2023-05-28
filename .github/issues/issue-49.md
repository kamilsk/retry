---
code: null
id: MDU6SXNzdWUyMjM2Mzg2NTM=
databaseId: 223638653
number: 49
url: https://github.com/kamilsk/retry/issues/49
title: "fix gometalinter"
issueType: null
labels:
  - "type:refactoring"
assignees:
  - kamilsk
milestone: "[[milestone-1]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-04-23T14:14:11Z
updatedAt: 2017-04-26T16:19:12Z
lastEditedAt: 2017-04-26T16:19:10Z
closedAt: 2017-04-26T16:19:12Z
---

# fix gometalinter

- [x] `cmd/retry/flag/flag.go:59:6:warning: type name will be used as flag.FlagSet by other packages, and that stutters; considercalling this Set (golint)`
- [x] `cmd/retry/go17_helper.go:54::error: the cancel function returned by context.WithTimeout should be called, not discarded, to avoid a context leak (vet)`
- [x] `cmd/retry/go17_helper.go:37:10:warning: error return value not checked (cl.Parse(os.Args[1:])) (errcheck)`
