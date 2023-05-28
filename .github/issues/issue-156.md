---
code: null
id: MDU6SXNzdWU2MjU4NTE0NzQ=
databaseId: 625851474
number: 156
url: https://github.com/kamilsk/retry/issues/156
title: "Wrong examples in the README"
issueType: null
labels:
  - "type:docs"
assignees:
  - kamilsk
milestone: "[[milestone-6]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2020-05-27T16:38:05Z
updatedAt: 2020-05-30T09:02:03Z
lastEditedAt: null
closedAt: 2020-05-30T09:02:02Z
---

# Wrong examples in the README

I'm struggling to follow the examples in the README because I can't relate them to the package actual implementation:

1. `breaker.BreakBySignal` as long as the other breakers don't implement the `strategy.Breaker` interface (the `Err()` method is missing). This is also a problem in the implementation, that is backing me from using this breaker

2. The function (second parameter of `retry.Do`) must accept a context as parameter. This is not reflected in the documentation

<!-- 2020-05-27T18:24Z https://github.com/kamilsk/retry/issues/156#issuecomment-634854412
Hi, @MarcoVigna17! Thanks for the report.

> The function (second parameter of retry.Do) must accept a context as parameter. This is not reflected in the documentation

Yes, sorry for that, v5 is still unstable (the latest version is v5.0.0-rc5), and I didn't sync the readme. I fix that soon.

About the first point: I have an issue with [an update of the interface](https://github.com/kamilsk/breaker/issues/8), and it is partially implemented.

So, I can suggest using:
1. retry/v4 together with breaker any version
2. retry/v5-rc together with breaker/v1.2-rc
-->

<!-- 2020-05-27T18:27Z https://github.com/kamilsk/retry/issues/156#issuecomment-634855814
I added the Err() method to the breakers [here](https://github.com/kamilsk/breaker/commit/1aac9aeb339081a8f43c49b466c96ade5d742cf1).
-->

<!-- 2020-05-28T08:18Z https://github.com/kamilsk/retry/issues/156#issuecomment-635192761
Thanks for the feedback @kamilsk.

I will try what you suggested today
-->

<!-- 2020-05-29T18:36Z https://github.com/kamilsk/retry/issues/156#issuecomment-636125284
@MarcoVigna17 I updated the readme and added information about `breaker` compatibility here
7a5261c
-->

<!-- 2020-05-30T09:02Z https://github.com/kamilsk/retry/issues/156#issuecomment-636302200
I close the issue but feel free to reopen it if needed.
-->
