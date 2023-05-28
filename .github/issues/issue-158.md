---
code: null
id: MDU6SXNzdWU3Nzc1NDI4MDU=
databaseId: 777542805
number: 158
url: https://github.com/kamilsk/retry/issues/158
title: "Preserve context values"
issueType: null
labels:
  - "type:task"
assignees:
  - kamilsk
milestone: "[[milestone-6]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2021-01-02T23:36:07Z
updatedAt: 2021-01-03T11:31:22Z
lastEditedAt: null
closedAt: 2021-01-03T11:15:51Z
---

# Preserve context values

`retry.Do` doesn't preserve context values when calling action. Here is an example:
```
ctx := context.WithValue(context.Background(), "asdf", 3)
action := func(ctx context.Context) error {
	log.Println(ctx.Value("asdf")) // Doesn't print 3
	return nil
}

retry.Do(ctx, action)
```
Does it make sense to change `retry.Do` to preserve context values?

<!-- 2021-01-03T08:53Z https://github.com/kamilsk/retry/issues/158#issuecomment-753586981
Hi! Yes, `retry.Do` doesn't preserve context values when calling action because it expects `strategy.Breaker`, not `context.Context`.

It is possible to do that:

```go
-	ctx, cancel := context.WithCancel(context.Background())
+	ctx, is := breaker.(context.Context)
+	if !is {
+		ctx = context.Background()
+	}
+	ctx, cancel := context.WithCancel(ctx)
```

I'm not a fan of the Value part of the context, but could you provide some use cases to make a decision.
-->

<!-- 2021-01-03T10:39Z https://github.com/kamilsk/retry/issues/158#issuecomment-753597581
Hi Kamil,

I am used to storing correlation id in context so it gets prepended to all log entries. For instance when I am building a REST API, I use request id as correlation id and I can then easily inspect all log entries produced by a request even if they are interleaved with other entries. That's how I found out `retry.Do` isn't preserving values - correlation id was missing from some log entries. There are other ways to do it of course, but projects I work on do it like that.

Database transaction handle also gets stored in context.
-->

<!-- 2021-01-03T10:46Z https://github.com/kamilsk/retry/issues/158#issuecomment-753598422
Thanks, Bojan! I'm preparing the [patch](#160) and will add some tests to it.

Today it will be ready, the upcoming release is v5.0.0-rc6.
-->

<!-- 2021-01-03T11:12Z https://github.com/kamilsk/retry/issues/158#issuecomment-753601147
Thank you very much!
-->

<!-- 2021-01-03T11:31Z https://github.com/kamilsk/retry/issues/158#issuecomment-753603157
https://github.com/kamilsk/retry/releases/tag/v5.0.0-rc6
-->
