package retry

import "context"

func convert(breaker Breaker) context.Context {
	ctx, is := breaker.(context.Context)
	if !is {
		ctx = lite{context.Background(), breaker}
	}
	return ctx
}

type lite struct {
	context.Context
	breaker Breaker
}

func (ctx lite) Done() <-chan struct{} { return ctx.breaker.Done() }
func (ctx lite) Err() error            { return ctx.breaker.Err() }
