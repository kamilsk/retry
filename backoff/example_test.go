package backoff_test

import (
	"fmt"
	"time"

	"github.com/kamilsk/retry/backoff"
)

func ExampleIncremental() {
	algorithm := backoff.Incremental(15*time.Millisecond, 10*time.Millisecond)

	for i := uint(1); i <= 5; i++ {
		duration := algorithm(i)

		fmt.Printf("#%d attempt: %s\n", i, duration)
		// Output:
		// #1 attempt: 25ms
		// #2 attempt: 35ms
		// #3 attempt: 45ms
		// #4 attempt: 55ms
		// #5 attempt: 65ms
	}
}

func ExampleLinear() {
	algorithm := backoff.Linear(15 * time.Millisecond)

	for i := uint(1); i <= 5; i++ {
		duration := algorithm(i)

		fmt.Printf("#%d attempt: %s\n", i, duration)
		// Output:
		// #1 attempt: 15ms
		// #2 attempt: 30ms
		// #3 attempt: 45ms
		// #4 attempt: 60ms
		// #5 attempt: 75ms
	}
}

func ExampleExponential() {
	algorithm := backoff.Exponential(15*time.Millisecond, 3)

	for i := uint(1); i <= 5; i++ {
		duration := algorithm(i)

		fmt.Printf("#%d attempt: %s\n", i, duration)
		// Output:
		// #1 attempt: 45ms
		// #2 attempt: 135ms
		// #3 attempt: 405ms
		// #4 attempt: 1.215s
		// #5 attempt: 3.645s
	}
}

func ExampleBinaryExponential() {
	algorithm := backoff.BinaryExponential(15 * time.Millisecond)

	for i := uint(1); i <= 5; i++ {
		duration := algorithm(i)

		fmt.Printf("#%d attempt: %s\n", i, duration)
		// Output:
		// #1 attempt: 30ms
		// #2 attempt: 60ms
		// #3 attempt: 120ms
		// #4 attempt: 240ms
		// #5 attempt: 480ms
	}
}

func ExampleFibonacci() {
	algorithm := backoff.Fibonacci(15 * time.Millisecond)

	for i := uint(1); i <= 5; i++ {
		duration := algorithm(i)

		fmt.Printf("#%d attempt: %s\n", i, duration)
		// Output:
		// #1 attempt: 15ms
		// #2 attempt: 15ms
		// #3 attempt: 30ms
		// #4 attempt: 45ms
		// #5 attempt: 75ms
	}
}
