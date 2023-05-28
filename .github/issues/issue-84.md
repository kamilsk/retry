---
code: null
id: MDU6SXNzdWUyNTU1ODc0MTM=
databaseId: 255587413
number: 84
url: https://github.com/kamilsk/retry/issues/84
title: "idea: retry builder"
issueType: null
labels:
  - "type:feature"
  - "status:research"
assignees:
  - kamilsk
milestone: "[[milestone-5]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-09-06T12:33:30Z
updatedAt: 2019-08-14T07:33:44Z
lastEditedAt: null
closedAt: 2019-02-24T15:41:34Z
---

# idea: retry builder

`Retrier.Before.Try.Otherwise.Finally`

<!-- 2017-09-08T09:02Z https://github.com/kamilsk/retry/issues/84#issuecomment-328046550
The concept:

```go
ctx := pipeline.Context()
err := retry.New().
	Before(retry.Retry(ctx.Done(), action.Login, login.Strategies...)).
	Try(action.DoSomething).
	Try(action.DoSomethingElse).
	Parallel().
	Try(transaction.Commit).
	Otherwise(transaction.Rollback).
	Finally(connection.Close).
	Do(ctx.Done())
```
-->

<!-- 2017-09-08T09:43Z https://github.com/kamilsk/retry/issues/84#issuecomment-328057371
```go
type Retrier struct {
	parallel                   []retry.Action
	sequence                   []retry.Action
	before, otherwise, finally func() error
}

func (r *Retrier) Parallel() *Retrier {
	r.parallel = append(r.parallel, r.sequence...)
	r.sequence = []retry.Action{}
	return r
}

func (r *Retrier) Try(action retry.Action) *Retrier {
	r.sequence = append(r.sequence, action)
	return r
}

func (r *Retrier) Do(deadline <-chan struct{}, strategies ...strategy.Strategy) error {
	ch := make(chan error, 1)
	go func() {
		err := <-ch
		if r.finally != nil {
			err := r.finally
			...
		}
		close(...)
	}()
	if r.before != nil {
		r.before()
	}
	for _, action := range r.parallel {
		go func() {
			ch <- retry.Retry(deadline, action, strategies...)
		}()
	}
	...

	return nil
}
```
-->

<!-- 2017-09-08T12:13Z https://github.com/kamilsk/retry/issues/84#issuecomment-328087037
Maybe is it better to use next semantic:
```go
Try(action.DoSomething, retry.Parallel(action.DoSomethingElse, action.DoSomethingMore)).
```
-->

<!-- 2017-09-08T12:25Z https://github.com/kamilsk/retry/issues/84#issuecomment-328089331
good function approach: https://www.calhoun.io/using-functional-options-instead-of-method-chaining-in-go/
-->

<!-- 2017-09-08T12:27Z https://github.com/kamilsk/retry/issues/84#issuecomment-328089835
```go
Retry(
  Before(Do).Try(Action, Parallel(Action, ...), ...).With(Strategies...).Catch(Do).Finally(Do)
)
```
-->

<!-- 2017-12-15T06:40Z https://github.com/kamilsk/retry/issues/84#issuecomment-351926690
investigate https://github.com/golang/sync/tree/master/errgroup
-->

<!-- 2017-12-19T19:57Z https://github.com/kamilsk/retry/issues/84#issuecomment-352869376
intestigate https://github.com/lazada/awg
-->

<!-- 2018-12-01T20:13Z https://github.com/kamilsk/retry/issues/84#issuecomment-443455029
`Repeat(Action).Until(...).Success(callback).Failure(callback)`
-->
