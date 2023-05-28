---
code: null
id: MDU6SXNzdWUyODg1NzUwMDY=
databaseId: 288575006
number: 109
url: https://github.com/kamilsk/retry/issues/109
title: "goroutine leak"
issueType: null
labels:
  - "type:critical"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2018-01-15T12:15:45Z
updatedAt: 2019-02-17T10:14:43Z
lastEditedAt: 2018-01-15T12:27:55Z
closedAt: 2019-02-17T10:14:43Z
---

# goroutine leak

- [ ] https://github.com/kamilsk/retry/blob/master/channel.go#L42-L46
- [ ] https://github.com/kamilsk/retry/blob/master/context.go#L11-L12
- [ ] https://github.com/kamilsk/retry/blob/master/channel.go#L52-L63 - channel + cancel, because `time.After` is bad
- [ ] https://github.com/kamilsk/retry/blob/master/channel.go#L35-L49 - the same as above, because we don't need to listen signal when work is done
- [ ] investigate more

<!-- 2018-05-04T09:05Z https://github.com/kamilsk/retry/issues/109#issuecomment-386544130
https://go101.org/article/concurrent-common-mistakes.html
-->
