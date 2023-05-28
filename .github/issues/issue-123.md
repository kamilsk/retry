---
code: null
id: MDU6SXNzdWUzNjc0MDg5NjE=
databaseId: 367408961
number: 123
url: https://github.com/kamilsk/retry/issues/123
title: "panic dead lock"
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
createdAt: 2018-10-06T00:59:31Z
updatedAt: 2018-10-09T12:07:08Z
lastEditedAt: null
closedAt: 2018-10-06T12:03:49Z
---

# panic dead lock

If action panic then the done channel is never closed and it hangs forever.

https://github.com/kamilsk/retry/blob/272c5373f544785efa4a1fd557aadfe505469695/retry.go#L42

<!-- 2018-10-06T12:03Z https://github.com/kamilsk/retry/issues/123#issuecomment-427568568
@Kuqd, thank you for reporting! fixed by 21779a3
-->

<!-- 2018-10-09T12:07Z https://github.com/kamilsk/retry/issues/123#issuecomment-428166835
Glad I could help !
-->
