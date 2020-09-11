package hrw

import (
	"fmt"
	"math"
	"math/bits"
	"sort"

	"github.com/nspcc-dev/hrw/normalizer"
)

func allSame(ns []uint64) bool {
	for i := range ns {
		if ns[i] != ns[0] {
			return false
		}
	}
	return true
}

func SortByWeightU64(norm normalizer.Uint64Norm, nodes, weights []uint64, hash uint64) []uint64 {
	for i := range weights {
		weights[i] = norm.Normalize(weights[i])
	}
	return SortByWeightU64Normalized(nodes, weights, hash)
}

func SortByWeightU64Normalized(nodes, weights []uint64, hash uint64) []uint64 {
	if len(nodes) != len(weights) {
		panic(fmt.Errorf("lengths don't match: %d vs %d", len(nodes), len(weights)))
	}

	// if all nodes have the same distance then sort uniformly
	if allSame(nodes) {
		Sort(nodes, hash)
	}

	ind := make([]int, len(nodes))
	dist := make([]uint64, len(nodes))
	for i := range nodes {
		ind[i] = i
		dist[i] = distance(nodes[i], hash)
	}

	sort.Slice(ind, func(i, j int) bool {
		ii, jj := ind[i], ind[j]
		// `maxUint64 - distance` makes the shorter distance more valuable
		// it is necessary for operation with normalized values
		di := math.MaxUint64 - dist[ii]
		dj := math.MaxUint64 - dist[jj]
		wiH, wiL := bits.Mul64(weights[ii], di)
		wjH, wjL := bits.Mul64(weights[jj], dj)
		return wiH > wjH || wiH == wjH && wiL > wjL // higher distance must be placed lower to be first
	})

	res := make([]uint64, len(nodes))
	for i := range res {
		res[i] = nodes[ind[i]]
	}
	return res
}
