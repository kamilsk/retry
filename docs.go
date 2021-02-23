// Package retry provides the most advanced interruptible mechanism
// to perform actions repetitively until successful.
//
// The retry based on https://github.com/Rican7/retry but fully reworked
// and focused on integration with the https://github.com/kamilsk/breaker
// and the built-in https://pkg.go.dev/context package.
package retry
