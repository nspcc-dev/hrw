// Package hrw implements Rendezvous hashing.
// https://en.wikipedia.org/wiki/Rendezvous_hashing.
package hrw

import (
	"sort"

	"github.com/twmb/murmur3"
	"golang.org/x/exp/constraints"
)

// Hashable is something that can be hashed.
type Hashable interface{ Hash() uint64 }

// HashableBytes implements Hashable interface over
// raw data. Use [WrapBytes] to instantiate a correct
// byte slice wrapper.
type HashableBytes []byte

func (h HashableBytes) Hash() uint64 {
	return murmur3.Sum64(h)
}

// WrapBytes creates [HashableBytes] that implements
// [Hashable] interface over a raw data.
// Can be used for [Sort] and [SortWeighted].
func WrapBytes(b []byte) HashableBytes {
	return b
}

// Sort defines and sorts the scores for the provided hashable
// entities against the provided hashable object (in its general
// sense).
// See [Hashable], [HashableBytes] and https://en.wikipedia.org/wiki/Rendezvous_hashing.
func Sort[V, P Hashable](vv []V, object P) {
	oHash := object.Hash()

	var s sliceToSort[V, uint64]
	s.s = vv
	s.distances = make([]uint64, len(vv))

	for i := range vv {
		s.distances[i] = distance(vv[i].Hash(), oHash)
	}

	sort.Stable(&s)
}

// SortWeighted is the same as [Sort] but allows using weights for
// a slice being sorted. A weight is applied to the corresponding
// element's score in the resulting slice. A weight allows modifying
// the default (equal to any element) probability of an element to
// win HRW sorting.
// Value slice's length and weight slice's length MUST be the same.
// Weights MUST be in [0.0; 1.0] range.
func SortWeighted[V, P Hashable, W constraints.Float](vv []V, weights []W, object P) {
	if len(vv) != len(weights) {
		return
	}

	if allSameF(weights) {
		Sort(vv, object)
		return
	}

	oHash := object.Hash()

	var s sliceToSort[V, W]
	s.s = vv
	s.distances = make([]W, len(vv))

	for i := range vv {
		// the distance is a bad characteristic in our case (we sort in ascending order)
		// so a bigger weight should lower the distance more
		s.distances[i] = W(distance(vv[i].Hash(), oHash)) / weights[i]
	}

	sort.Stable(&s)
}

type distancesValue interface {
	constraints.Unsigned | constraints.Float
}

type sliceToSort[V any, W distancesValue] struct {
	s         []V
	distances []W
}

func (s *sliceToSort[V, _]) Len() int {
	return len(s.s)
}

func (s *sliceToSort[V, _]) Less(i, j int) bool {
	return s.distances[i] < s.distances[j]
}

func (s *sliceToSort[V, _]) Swap(i, j int) {
	s.s[i], s.s[j] = s.s[j], s.s[i]
	s.distances[i], s.distances[j] = s.distances[j], s.distances[i]
}

func distance(x uint64, y uint64) uint64 {
	acc := x ^ y
	// here used mmh3 64 bit finalizer
	// https://github.com/aappleby/smhasher/blob/61a0530f28277f2e850bfc39600ce61d02b518de/src/MurmurHash3.cpp#L81
	acc ^= acc >> 33
	acc = acc * 0xff51afd7ed558ccd
	acc ^= acc >> 33
	acc = acc * 0xc4ceb9fe1a85ec53
	acc ^= acc >> 33
	return acc
}

func allSameF[W constraints.Float](fs []W) bool {
	for i := range fs {
		if fs[i] != fs[0] {
			return false
		}
	}
	return true
}
