package normalizer

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

// In this test we check that result of uint64 normalizer
// corresponds to that of FloatNorm one.
func TestUint64Normalizer(t *testing.T) {
	numbers := []uint64{1, 10, 100, 1000, 1000_000, 1000_000_000_000}
	testCases := []struct {
		name  string
		count int
		u     func(uint64) Uint64Norm
		f     func(float64) FloatNorm
	}{
		{"LogRatio", 11, NewLogRatioU64, NewLogRatioF64},
		{"Max", 11, NewMaxU64, NewMaxF64},
		{"ReverseMin", -11, NewReverseMinU64, NewReverseMinF64},
		{"Sigmoid", 11, NewSigmoidU64, NewSigmoidF64},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, num := range numbers {
				uNorm := tc.u(num)
				fNorm := tc.f(float64(num))
				steps := getSteps(num, tc.count)
				for _, v := range steps {
					uRes := uNorm.Normalize(v)
					fRes := fNorm.Normalize(float64(v))
					require.InDelta(t, fRes, ToFloat(uRes), 0.000_001,
						"value: %d, norm: %d", v, num)
				}
			}
		})
	}
}

func getSteps(max uint64, count int) []uint64 {
	steps := []uint64{max}
	if count > 0 {
		factor := (max / uint64(count))
		for i := 0; i < count; i++ {
			if steps[len(steps)-1] != uint64(i)*factor {
				steps = append(steps, uint64(i)*factor)
			}
		}
	} else {
		factor := (uint64(math.MaxUint64)-max)/uint64(count) + 1
		for i := 0; i < -count; i++ {
			val := max + uint64(i)*factor
			if steps[len(steps)-1] != val {
				steps = append(steps, val)
			}
		}
	}
	return steps
}

func TestConstU64_Normalize(t *testing.T) {
	norm := NewConstU64(123)
	require.Equal(t, uint64(123), norm.Normalize(10))
	require.Equal(t, uint64(123), norm.Normalize(1000))
}
