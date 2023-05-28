---
code: null
id: MDU6SXNzdWU4MTMxODIyMzM=
databaseId: 813182233
number: 171
url: https://github.com/kamilsk/retry/issues/171
title: "wat moment: Schrödinger's bug"
issueType: null
labels:
  - "type:critical"
  - "status:research"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-6]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2021-02-22T06:41:05Z
updatedAt: 2021-02-22T15:44:49Z
lastEditedAt: 2021-02-22T15:42:59Z
closedAt: 2021-02-22T15:44:48Z
---

# wat moment: Schrödinger's bug

I ran into an unpleasant problem:

```go
ctx = // some top-level context

client = // some http client
req := http.NewRequestWithContext(...)

var resp *http.Response

action := func(ctx context.Context) error {
    var err error
    resp, err = client.Do(req.WithContext(ctx)) // it's important to wrap action's ctx to try to reproduce the problem
    return err
}

if err := retry.Do(ctx, action, ...); err != nil {
    return nil, err
}
defer safe.Close(io.Discarded(resp.Body), unsafe.Ignore)

var result map[string]interface{}
if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
    return nil, err
}
return result, nil // it's fluky because I also may have a "context canceled" error
```

It's happening because I cancel the action's context on success https://github.com/kamilsk/retry/blob/13994452f848aa8b4f15f17cb82c55703f465b1b/retry.go#L69 and https://github.com/golang/go/blob/0f66fb7b856b02497cf801ce72d80f375f53358b/src/net/http/client.go#L928-L960

I can't don't do it cause of goroutine leaking

What I have to do:

- [x] add a workaround to avoid it now
- [x] define a valid solution
