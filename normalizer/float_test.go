package normalizer

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

const eps = 0.000_001

func TestSigmoidNorm_Normalize(t *testing.T) {
	t.Run("sigmoid norm must equal to 1/2 at `scale`", func(t *testing.T) {
		norm := NewSigmoidF64(1)
		require.InEpsilon(t, 0.5, norm.Normalize(1), eps)

		norm = NewSigmoidF64(10)
		require.InEpsilon(t, 0.5, norm.Normalize(10), eps)
	})

	t.Run("sigmoid norm must be less than 1", func(t *testing.T) {
		norm := NewSigmoidF64(2)
		require.True(t, norm.Normalize(100) < 1)
		require.True(t, norm.Normalize(math.MaxFloat64) <= 1)
	})

	t.Run("sigmoid norm must be monotonic", func(t *testing.T) {
		norm := NewSigmoidF64(5)
		for i := 0; i < 5; i++ {
			a, b := rand.Float64(), rand.Float64()
			if b < a {
				a, b = b, a
			}
			require.True(t, norm.Normalize(a) <= norm.Normalize(b))
		}
	})
}

func TestReverseMinF64_Normalize(t *testing.T) {
	t.Run("reverseMin norm should not panic", func(t *testing.T) {
		norm := NewReverseMinF64(0)
		require.NotPanics(t, func() { norm.Normalize(0) })

		norm = NewReverseMinF64(1)
		require.NotPanics(t, func() { norm.Normalize(0) })
	})

	t.Run("reverseMin norm should equal 1 at min value", func(t *testing.T) {
		norm := NewReverseMinF64(10)
		require.InEpsilon(t, 1.0, norm.Normalize(10), eps)
	})
}

func TestMaxF64_Normalize(t *testing.T) {
	t.Run("max norm should not panic", func(t *testing.T) {
		norm := NewMaxF64(0)
		require.NotPanics(t, func() { norm.Normalize(1) })

		norm = NewMaxF64(1)
		require.NotPanics(t, func() { norm.Normalize(0) })
	})

	t.Run("max norm should equal 1 at max value", func(t *testing.T) {
		norm := NewMaxF64(10)
		require.InEpsilon(t, 1.0, norm.Normalize(10), eps)
	})
}

func TestConstF64_Normalize(t *testing.T) {
	norm := NewConstF64(123)
	require.Equal(t, float64(123), norm.Normalize(10))
	require.Equal(t, float64(123), norm.Normalize(1000))
}
