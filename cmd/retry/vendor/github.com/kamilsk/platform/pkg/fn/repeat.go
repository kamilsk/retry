package fn

import "github.com/kamilsk/platform/pkg/math"

// Repeat repeats the action the required number of times.
//
//  func FillByValue(slice []int, value, count int) []int {
//  	Repeat(func () { slice = append(slice, value) }, count)
//  	return slice
//  }
//
func Repeat(action func(), times int) {
	for range math.Sequence(times) {
		action()
	}
}
