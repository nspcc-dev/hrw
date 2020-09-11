package hrw

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/nspcc-dev/hrw/normalizer"
	"github.com/stretchr/testify/require"
)

func copyUint64s(us []uint64) []uint64 {
	res := make([]uint64, len(us))
	copy(res, us)
	return res
}

func TestSortByWeightU64(t *testing.T) {
	const (
		size    = 10
		keys    = 100000
		percent = 0.03
	)

	var (
		i            uint64
		a, w, result [size]uint64
		key          = make([]byte, 16)
	)

	for i = 0; i < size; i++ {
		a[i] = i
		w[int(i)] = 10
	}
	w[0] = 100
	norm := normalizer.NewLogRatioU64(w[0])
	for i = 0; i < keys; i++ {
		binary.BigEndian.PutUint64(key, i+size)
		hash := Hash(key)
		ns := SortByWeightU64(norm, a[:], copyUint64s(w[:]), hash)
		for j := range ns {
			result[ns[j]] += uint64(len(ns) - j)
		}
	}

	cutResult := result[1:]
	var total uint64
	for i := range cutResult {
		total += cutResult[i]
	}

	var chi2 float64
	mean := float64(total) / float64(len(cutResult))
	delta := mean * percent
	for node, count := range cutResult {
		d := mean - float64(count)
		chi2 += math.Pow(float64(count)-mean, 2) / mean
		require.True(t, d < delta && (0-d) < delta,
			"Node %d received %d keys, expected %.0f (+/- %.2f)", node, count, mean, delta)
	}
	require.True(t, chi2 < chiTable[size-1],
		"Chi2 condition for .9 is not met (expected %.2f <= %.2f)", chi2, chiTable[size-1])
}
