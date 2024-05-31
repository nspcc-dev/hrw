package normalizer

import (
	"math"
	"math/bits"
)

// ToFloat transforms u from [0..math.MaxUint64] to a float in [0..1].
func ToFloat(u uint64) float64 {
	return float64(u) / float64(math.MaxUint64)
}

// Div returns a / b mapped from [0..1] to [0..math.MaxUint64].
func Div(a, b uint64) uint64 {
	hi, lo := bits.Mul64(a, math.MaxUint64)
	if hi > b {
		return math.MaxUint64
	}
	q, _ := bits.Div64(hi, lo, b)
	return q
}

// logFactor is used to gain some precision when mapping from logarithm
// to uint64 domain.
const logFactor = 2 << 16

func log2(u uint64) uint64 {
	return uint64(math.Log2(float64(u)) * logFactor)
}
