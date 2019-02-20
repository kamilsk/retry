package math

import (
	"math"
	"sort"
)

// Sequence returns an empty slice with the specified size.
//
//  for range Sequence(5) {
//  	// do something five times
//  }
//
func Sequence(size int) []struct{} {
	return make([]struct{}, size)
}

// BatchSequence returns a specific slice for batch iteration on another slice.
//
//  batch, target := 100, make([]int, 1024)
//  for step, end := range BatchSequence(len(target), batch) {
//  	process(target[batch*step:end])
//  }
//
func BatchSequence(size, batch int) []int {
	count := size / batch
	if size%batch != 0 {
		count++
	}
	batches := make([]int, count)
	for i := 0; i < count; i++ {
		border := (i + 1) * batch
		if border > size {
			border = size
		}
		batches[i] = border
	}
	return batches
}

// Reduce wraps sequence to perform some aggregate operations above it.
//
//  func Acquire(places ...int) {
//  	for range Sequence(Reduce(places...).Sum()) {
//  		semaphore <- struct{}{}
//  	}
//  }
//
func Reduce(sequence ...int) interface {
	// Average returns an average value of the sequence.
	Average() float64
	// Count returns the sequence length.
	Count() int
	// Maximum returns a maximum value of the sequence.
	Maximum() int
	// Median returns a median value of the sequence.
	Median() float64
	// Minimum returns a minimum value of the sequence.
	Minimum() int
	// Sum returns a sum of the sequence.
	Sum() int
} {
	return reducer(sequence)
}

type reducer []int

// Average returns an average value of the sequence.
func (sequence reducer) Average() float64 {
	if len(sequence) == 0 {
		return 0
	}
	return float64(sequence.Sum()) / float64(len(sequence))
}

// Count returns the sequence length.
func (sequence reducer) Count() int {
	return len(sequence)
}

// Maximum returns a maximum value of the sequence.
func (sequence reducer) Maximum() int {
	if len(sequence) == 0 {
		return 0
	}
	max := math.MinInt64
	for _, num := range sequence {
		if num > max {
			max = num
		}
	}
	return max
}

// Median returns a median value of the sequence.
func (sequence reducer) Median() float64 {
	size := len(sequence)
	if size == 0 {
		return 0
	}
	sorted := append(make([]int, 0, size), sequence...)
	sort.Ints(sorted)
	if size%2 == 0 {
		return (float64(sorted[size/2-1]) + float64(sorted[size/2])) / 2
	}
	return float64(sorted[size/2])
}

// Minimum returns a minimum value of the sequence.
func (sequence reducer) Minimum() int {
	if len(sequence) == 0 {
		return 0
	}
	min := math.MaxInt64
	for _, num := range sequence {
		if num < min {
			min = num
		}
	}
	return min
}

// Sum returns a sum of the sequence.
func (sequence reducer) Sum() int {
	sum := 0
	for _, num := range sequence {
		sum += num
	}
	return sum
}
