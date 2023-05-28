---
code: null
id: MDU6SXNzdWUyNDUxMjM0OTk=
databaseId: 245123499
number: 77
url: https://github.com/kamilsk/retry/issues/77
title: "handle os.Interrupt signal"
issueType: null
labels:
  - "type:task"
assignees:
  - kamilsk
milestone: "[[milestone-3]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-07-24T15:45:15Z
updatedAt: 2017-11-30T11:44:10Z
lastEditedAt: null
closedAt: 2017-11-30T11:44:10Z
---

# handle os.Interrupt signal

```
c := make(chan os.Signal)
signal.Notify(c, os.Interrupt)
```

<!-- 2017-07-31T12:45Z https://github.com/kamilsk/retry/issues/77#issuecomment-319056696
provide a useful abstraction for that. use sync.Once to initialize it
-->

<!-- 2017-11-30T11:44Z https://github.com/kamilsk/retry/issues/77#issuecomment-348165203
done by e7612610f27c0babd097e568a9673bf11ce29c82
-->
