package safecast

import (
	"math"
)

// IntToInt32 clamps an int to the valid range of int32.
func IntToInt32(number int) int32 {
	if number > math.MaxInt32 {
		return math.MaxInt32
	}

	if number < math.MinInt32 {
		return math.MinInt32
	}

	return int32(number)
}
