package normalizer

import (
	"math"
)

type (
	// Uint64Norm transforms uint64 weight to a some kind of normalized value.
	Uint64Norm interface {
		Normalize(w uint64) uint64
	}

	reverseMinU64 struct {
		min uint64
	}

	maxU64 struct {
		max uint64
	}

	sigmoidU64 struct {
		scale uint64
	}

	constU64 struct {
		value uint64
	}

	logRatioU64 struct {
		maxLog uint64
	}
)

var (
	_ Uint64Norm = (*constU64)(nil)
	_ Uint64Norm = (*logRatioU64)(nil)
	_ Uint64Norm = (*maxU64)(nil)
	_ Uint64Norm = (*reverseMinU64)(nil)
	_ Uint64Norm = (*sigmoidU64)(nil)
)

// NewReverseMinU64 returns a normalizer which
// normalize values in range of 0.0 to 1.0 to a minimum value.
func NewReverseMinU64(min uint64) Uint64Norm {
	return &reverseMinU64{min: min}
}

// Normalize implements Uint64Norm.
func (r *reverseMinU64) Normalize(w uint64) uint64 {
	if w == 0 {
		return 0
	}
	return Div(r.min, w)
}

// NewMaxU64 returns a normalizer which
// normalize values in range of 0.0 to 1.0 to a maximum value.
func NewMaxU64(max uint64) Uint64Norm {
	return &maxU64{max: max}
}

// Normalize implements Uint64Norm.
func (r *maxU64) Normalize(w uint64) uint64 {
	if r.max == 0 {
		return 0
	}
	return Div(w, r.max)
}

// NewSigmoidU64 returns a normalizer which
// normalize values in range of 0.0 to 1.0 to a scaled sigmoid.
func NewSigmoidU64(scale uint64) Uint64Norm {
	if scale == 0 {
		panic("zero scale")
	}
	return &sigmoidU64{scale: scale}
}

// Normalize implements Uint64Norm.
func (r *sigmoidU64) Normalize(w uint64) uint64 {
	x := Div(w, r.scale)
	// 1/(x+1) = 1 - 1/(x+1) corresponds to math.MaxUint64 - math.MaxUint64 / (x+math.MaxUint64)
	return math.MaxUint64 - Div(math.MaxUint64/2, x/2+math.MaxUint64/2)
}

// NewConstU64 returns a normalizer which
// returns a constant values
func NewConstU64(value uint64) Uint64Norm {
	return &constU64{value: value}
}

// Normalize implements Uint64Norm.
func (r *constU64) Normalize(_ uint64) uint64 {
	return r.value
}

// NewLogRatio returns a normalizer for which norm
// is the ratio of value to log2(max).
func NewLogRatioU64(max uint64) Uint64Norm {
	return &logRatioU64{maxLog: log2(max)}
}

// Normalize implements Uint64Norm.
func (r *logRatioU64) Normalize(x uint64) uint64 {
	if r.maxLog == 0 || x == 0 {
		return 0
	}
	return Div(log2(x), r.maxLog)
}
