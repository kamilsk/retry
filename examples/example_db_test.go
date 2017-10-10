// +build go1.8

package examples

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/kamilsk/retry"
	"github.com/kamilsk/retry/backoff"
	"github.com/kamilsk/retry/jitter"
	"github.com/kamilsk/retry/strategy"
)

type dummyDriver struct {
	conn *dummyConn
}

func (d *dummyDriver) Open(name string) (driver.Conn, error) {
	return d.conn, nil
}

type dummyConn struct {
	counter int
	ping    chan error
}

func (c *dummyConn) Prepare(string) (driver.Stmt, error) { return nil, nil }

func (c *dummyConn) Close() error { return nil }

func (c *dummyConn) Begin() (driver.Tx, error) { return nil, nil }

func (c *dummyConn) Ping(context.Context) error {
	c.counter++
	return <-c.ping
}

// This example shows how to use the library to restore database connection.
func Example_dbConnectionRestore() {
	d := &dummyDriver{conn: &dummyConn{ping: make(chan error, 10)}}
	for i := 0; i < cap(d.conn.ping); i++ {
		d.conn.ping <- nil
	}
	sql.Register("sqlite", d)

	shutdown := make(chan struct{})

	MustOpen := func() *sql.DB {
		db, err := sql.Open("sqlite", "./sqlite.db")
		if err != nil {
			panic(err)
		}
		return db
	}

	go func(db *sql.DB, ctx context.Context, shutdown chan<- struct{}, attempt uint, frequency time.Duration) {
		defer func() {
			if r := recover(); r != nil {
				shutdown <- struct{}{}
			}
		}()

		ping := func(attempt uint) error {
			return db.Ping()
		}
		strategies := []strategy.Strategy{
			strategy.Limit(attempt),
			strategy.BackoffWithJitter(
				backoff.Incremental(time.Millisecond, time.Millisecond),
				jitter.NormalDistribution(rand.New(rand.NewSource(time.Now().UnixNano())), 2.0),
			),
		}

		for {
			if err := retry.Retry(ctx.Done(), ping, strategies...); err != nil {
				panic(err)
			}
			time.Sleep(frequency)
		}
	}(MustOpen(), context.Background(), shutdown, 1, time.Millisecond)

	d.conn.ping <- errors.New("done")
	<-shutdown

	fmt.Printf("number of ping calls: %d", d.conn.counter)
	// Output: number of ping calls: 11
}
