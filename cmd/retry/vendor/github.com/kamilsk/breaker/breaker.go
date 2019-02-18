package breaker // import "github.com/kamilsk/breaker"

import (
	"context"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

// Interface carries a cancellation signal to break an action execution.
//
// Example based on github.com/kamilsk/retry package:
//
//  if err := retry.Retry(breaker.BreakByTimeout(time.Minute), action); err != nil {
//  	log.Fatal(err)
//  }
//
// Example based on github.com/kamilsk/semaphore package:
//
//  if err := semaphore.Acquire(breaker.BreakByTimeout(time.Minute), 5); err != nil {
//  	log.Fatal(err)
//  }
//
type Interface interface {
	// Done returns a channel that's closed when a cancellation signal occurred.
	Done() <-chan struct{}
	// Close closes the Done channel and releases resources associated with it.
	Close()
	// trigger is a private method to guarantee that the Breakers come from
	// this package and all of them return a valid Done channel.
	trigger() Interface
}

// BreakByDeadline closes the Done channel when the deadline occurs.
func BreakByDeadline(deadline time.Time) Interface {
	timeout := time.Until(deadline)
	if timeout < 0 {
		return closedBreaker()
	}
	return newTimedBreaker(timeout).trigger()
}

// BreakBySignal closes the Done channel when signals will be received.
func BreakBySignal(sig ...os.Signal) Interface {
	if len(sig) == 0 {
		return closedBreaker()
	}
	return newSignaledBreaker(sig).trigger()
}

// BreakByTimeout closes the Done channel when the timeout happens.
func BreakByTimeout(timeout time.Duration) Interface {
	if timeout < 0 {
		return closedBreaker()
	}
	return newTimedBreaker(timeout).trigger()
}

// Multiplex combines multiple Breakers into one.
func Multiplex(breakers ...Interface) Interface {
	if len(breakers) == 0 {
		return closedBreaker()
	}
	return newMultiplexedBreaker(breakers).trigger()
}

// MultiplexTwo combines two Breakers into one.
// This is the optimized version of more generic Multiplex.
func MultiplexTwo(one, two Interface) Interface {
	br := newBreaker()
	go func() {
		defer br.Close()
		select {
		case <-one.Done():
		case <-two.Done():
		}
	}()
	return br
}

// MultiplexThree combines three Breakers into one.
// This is the optimized version of more generic Multiplex.
func MultiplexThree(one, two, three Interface) Interface {
	br := newBreaker()
	go func() {
		defer br.Close()
		select {
		case <-one.Done():
		case <-two.Done():
		case <-three.Done():
		}
	}()
	return br
}

// WithContext returns a new Breaker and an associated Context derived from ctx.
func WithContext(ctx context.Context) (Interface, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	br := &contextBreaker{cancel, ctx.Done()}
	return br.trigger(), ctx
}

func closedBreaker() *breaker {
	br := newBreaker()
	br.Close()
	return br
}

func newBreaker() *breaker {
	return &breaker{signal: make(chan struct{})}
}

type breaker struct {
	closer   sync.Once
	signal   chan struct{}
	released int32
}

// Done returns a channel that's closed when a cancellation signal occurred.
func (br *breaker) Done() <-chan struct{} {
	return br.signal
}

// Close closes the Done channel and releases resources associated with it.
func (br *breaker) Close() {
	br.closer.Do(func() { close(br.signal) })
}

func (br *breaker) trigger() Interface {
	return br
}

type contextBreaker struct {
	cancel context.CancelFunc
	signal <-chan struct{}
}

// Done returns a channel that's closed when a cancellation signal occurred.
func (br *contextBreaker) Done() <-chan struct{} {
	return br.signal
}

// Close closes the Done channel and releases resources associated with it.
func (br *contextBreaker) Close() {
	br.cancel()
}

func (br *contextBreaker) trigger() Interface {
	return br
}

func newMultiplexedBreaker(entries []Interface) Interface {
	return &multiplexedBreaker{newBreaker(), entries}
}

type multiplexedBreaker struct {
	*breaker
	entries []Interface
}

// Close closes the Done channel and releases resources associated with it.
func (br *multiplexedBreaker) Close() {
	br.closer.Do(func() {
		each(br.entries).Close()
		close(br.signal)
	})
}

// trigger starts listening all Done channels of multiplexed Breakers.
func (br *multiplexedBreaker) trigger() Interface {
	go func() {
		brs := make([]reflect.SelectCase, 0, len(br.entries))
		for _, br := range br.entries {
			brs = append(brs, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(br.Done())})
		}
		reflect.Select(brs)
		br.Close()
		atomic.StoreInt32(&br.released, 1)
	}()
	return br
}

func newSignaledBreaker(signals []os.Signal) Interface {
	return &signaledBreaker{newBreaker(), make(chan os.Signal, len(signals)), signals}
}

type signaledBreaker struct {
	*breaker
	relay   chan os.Signal
	signals []os.Signal
}

// Close closes the Done channel and releases resources associated with it.
func (br *signaledBreaker) Close() {
	br.closer.Do(func() {
		signal.Stop(br.relay)
		close(br.signal)
	})
}

// trigger starts listening required signals to close the Done channel.
func (br *signaledBreaker) trigger() Interface {
	go func() {
		signal.Notify(br.relay, br.signals...)
		select {
		case <-br.relay:
		case <-br.signal:
		}
		br.Close()
		atomic.StoreInt32(&br.released, 1)
	}()
	return br
}

func newTimedBreaker(timeout time.Duration) Interface {
	return &timedBreaker{newBreaker(), time.NewTimer(timeout)}
}

type timedBreaker struct {
	*breaker
	*time.Timer
}

// Close closes the Done channel and releases resources associated with it.
func (br *timedBreaker) Close() {
	br.closer.Do(func() {
		br.Timer.Stop()
		close(br.signal)
	})
}

// trigger starts listening internal timer to close the Done channel.
func (br *timedBreaker) trigger() Interface {
	go func() {
		select {
		case <-br.Timer.C:
		case <-br.signal:
		}
		br.Close()
		atomic.StoreInt32(&br.released, 1)
	}()
	return br
}

type each []Interface

// Close closes all Done channels of a list of Breakers
// and releases resources associated with them.
func (list each) Close() {
	for _, br := range list {
		br.Close()
	}
}
