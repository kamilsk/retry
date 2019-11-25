// +build draft

package internal

import "github.com/kamilsk/retry/v4"

// Interface defines a behavior of stateful executor of Actions in parallel.
type Interface interface {
	Try(retry.Breaker, retry.Action, ...retry.How) Interface
}
