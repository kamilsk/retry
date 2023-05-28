---
code: null
id: MDU6SXNzdWU0ODA1OTg2MzI=
databaseId: 480598632
number: 136
url: https://github.com/kamilsk/retry/issues/136
title: "concept: flow, and add more deprecated marks"
issueType: null
labels:
  - "type:feature"
  - "status:pending"
assignees:
  - kamilsk
milestone: "[[milestone-5]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2019-08-14T10:20:43Z
updatedAt: 2020-02-10T15:36:27Z
lastEditedAt: null
closedAt: 2020-02-10T15:36:27Z
---

# concept: flow, and add more deprecated marks

Flow: `retry.What(action).How(strategies).Do(error handlers ?)`.

Error handler: `func(error) error`, for example, classifier or logger.

<!-- 2020-01-08T06:40Z https://github.com/kamilsk/retry/issues/136#issuecomment-571912870
`try.Do(...).Catch(Matching).Catch(Matching).Finally(...)`
-->

<!-- 2020-02-10T15:36Z https://github.com/kamilsk/retry/issues/136#issuecomment-584181920
won't do
-->
