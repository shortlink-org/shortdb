package pkg

import (
	"math"
)

// SafeIntToInt32 clamps an int to the valid range of int32.
func SafeIntToInt32(n int) int32 {
	if n > math.MaxInt32 {
		return math.MaxInt32
	}

	if n < math.MinInt32 {
		return math.MinInt32
	}

	return int32(n)
}
