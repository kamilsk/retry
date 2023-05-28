---
code: null
id: MDU6SXNzdWU4MTIyMzU4NDc=
databaseId: 812235847
number: 169
url: https://github.com/kamilsk/retry/issues/169
title: "unexpected Limit behavior"
issueType: null
labels:
  - "type:bug"
  - "status:research"
  - "help wanted"
  - "difficulty: easy"
assignees:
  - kamilsk
milestone: "[[milestone-6]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2021-02-19T18:21:37Z
updatedAt: 2021-02-21T11:54:26Z
lastEditedAt: null
closedAt: 2021-02-21T11:54:26Z
---

# unexpected Limit behavior

origin: https://github.com/Rican7/retry/blob/272ad122d6e5ce1be757544007cf8bcd1c9c9ab0/strategy/strategy.go#L25-L29

what if:

```go
var resp *http.Response

action := func(...) error { fill resp }

if err := retry.Do(action, strategy.Limit(0)); err != nil {
    return err
}

json.NewEncoder(resp.Body).Encode() <- nil pointer, because resp is nil, action didn't run
```

<!-- 2021-02-19T18:27Z https://github.com/kamilsk/retry/issues/169#issuecomment-782256222
user can define custom Strategy and the problem will repeat. as a solution I have to define default err = was not called
-->
