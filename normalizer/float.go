package normalizer

import "math"

type (
	// FloatNorm transforms float64 weight to a some kind of normalized value.
	FloatNorm interface {
		Normalize(w float64) float64
	}

	reverseMinF64 struct {
		min float64
	}

	maxF64 struct {
		max float64
	}

	sigmoidF64 struct {
		scale float64
	}

	constF64 struct {
		value float64
	}

	logRatioF64 struct {
		maxLog float64
	}
)

var (
	_ FloatNorm = (*constF64)(nil)
	_ FloatNorm = (*logRatioF64)(nil)
	_ FloatNorm = (*maxF64)(nil)
	_ FloatNorm = (*reverseMinF64)(nil)
	_ FloatNorm = (*sigmoidF64)(nil)
)

// NewReverseMinF64 returns a normalizer which
// normalize values in range of 0.0 to 1.0 to a minimum value.
func NewReverseMinF64(min float64) FloatNorm {
	return &reverseMinF64{min: min}
}

func (r *reverseMinF64) Normalize(w float64) float64 {
	if w == 0 {
		return 0
	}
	return r.min / w
}

// NewMaxF64 returns a normalizer which
// normalize values in range of 0.0 to 1.0 to a maximum value.
func NewMaxF64(max float64) FloatNorm {
	return &maxF64{max: max}
}

func (r *maxF64) Normalize(w float64) float64 {
	if r.max == 0 {
		return 0
	}
	return w / r.max
}

// NewSigmoidF64 returns a normalizer which
// normalize values in range of 0.0 to 1.0 to a scaled sigmoid.
func NewSigmoidF64(scale float64) FloatNorm {
	if scale == 0 {
		panic("zero scale")
	}
	return &sigmoidF64{scale: scale}
}

func (r *sigmoidF64) Normalize(w float64) float64 {
	x := w / r.scale
	return x / (1 + x)
}

// NewConstF64 returns a normalizer which
// returns a constant values
func NewConstF64(value float64) FloatNorm {
	return &constF64{value: value}
}

func (r *constF64) Normalize(_ float64) float64 {
	return r.value
}

// NewLogRatio returns a normalizer for which norm
// is the ratio of value to log2(max).
func NewLogRatioF64(max float64) FloatNorm {
	return &logRatioF64{maxLog: math.Log2(max)}
}

func (r *logRatioF64) Normalize(x float64) float64 {
	if r.maxLog == 0 || x == 0 {
		return 0
	}
	return math.Log2(x) / r.maxLog
}
