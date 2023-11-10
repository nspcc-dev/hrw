package hrw

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Compatibility tests to keep behavior the same
// compared with the v1-versioned library functions.

func TestSortSliceByIndex(t *testing.T) {
	actual := []string{"a", "b", "c", "d", "e", "f"}
	expect := []string{"e", "a", "c", "f", "d", "b"}

	testSlice := stringsToHashAndValueWithIndex(actual)
	Sort(testSlice, WrapBytes(testKey))
	result := hashAndValueToStrings(testSlice)

	require.Equal(t, expect, result)
}

func TestSortSliceByWeightIndex(t *testing.T) {
	actual := []string{"a", "b", "c", "d", "e", "f"}
	weights := []float64{1, 1, 1, 0.2, 0.2, 0.2}
	expect := []string{"a", "c", "b", "e", "f", "d"}

	testSlice := stringsToHashAndValueWithIndex(actual)
	SortWeighted(testSlice, weights, WrapBytes(testKey))
	result := hashAndValueToStrings(testSlice)

	require.Equal(t, expect, result)
}

func TestSortSliceByValue(t *testing.T) {
	actual := []string{"a", "b", "c", "d", "e", "f"}
	expect := []string{"d", "f", "c", "b", "a", "e"}

	testSlice := stringsToHashAndValueWithValue(actual)
	Sort(testSlice, WrapBytes(testKey))
	result := hashAndValueToStrings(testSlice)

	require.Equal(t, expect, result)
}

type hashAndValue struct {
	hash uint64
	val  string
}

func (h hashAndValue) Hash() uint64 {
	return h.hash
}

func stringsToHashAndValueWithIndex(ss []string) []hashAndValue {
	res := make([]hashAndValue, 0, len(ss))
	for i, s := range ss {
		res = append(res, hashAndValue{hash: uint64(i), val: s})
	}

	return res
}

func stringsToHashAndValueWithValue(ss []string) []hashAndValue {
	res := make([]hashAndValue, 0, len(ss))
	for _, s := range ss {
		res = append(res, hashAndValue{hash: WrapBytes([]byte(s)).Hash(), val: s})
	}

	return res
}

func hashAndValueToStrings(ss []hashAndValue) []string {
	res := make([]string, 0, len(ss))
	for _, s := range ss {
		res = append(res, s.val)
	}

	return res
}
