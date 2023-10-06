package hrw

import (
	"encoding/binary"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

type hashString string

func (h hashString) Hash() uint64 {
	return WrapBytes([]byte(h)).Hash()
}

var testKey = []byte("0xff51afd7ed558ccd")

func Example() {
	// given a set of servers
	servers := []Hashable{
		hashString("one.example.com"),
		hashString("two.example.com"),
		hashString("three.example.com"),
		hashString("four.example.com"),
		hashString("five.example.com"),
		hashString("six.example.com"),
	}

	// HRW can consistently select a uniformly-distributed set of servers for
	// any given key
	var key = []byte("/examples/object-key")

	Sort(servers, WrapBytes(key))
	for id := range servers {
		fmt.Printf("trying GET %s%s\n", servers[id], string(key))
	}

	// Output:
	// trying GET three.example.com/examples/object-key
	// trying GET two.example.com/examples/object-key
	// trying GET five.example.com/examples/object-key
	// trying GET six.example.com/examples/object-key
	// trying GET one.example.com/examples/object-key
	// trying GET four.example.com/examples/object-key
}

type hashableUint64 uint64

func (h hashableUint64) Hash() uint64 {
	return uint64(h)
}

func wrapUint64(uu []uint64) []Hashable {
	res := make([]Hashable, 0, len(uu))
	for _, u := range uu {
		res = append(res, hashableUint64(u))
	}

	return res
}

func TestSort(t *testing.T) {
	nodes := wrapUint64([]uint64{1, 2, 3, 4, 5})
	expected := wrapUint64([]uint64{4, 2, 5, 3, 1})

	Sort(nodes, WrapBytes(testKey))
	require.Equal(t, expected, nodes)
}

func TestDistribution(t *testing.T) {
	const (
		size    = 10
		keys    = 100000
		percent = 0.03
	)
	// We use χ2 method to determine similarity of distribution with uniform distribution.
	// χ2 = Σ((n-N)**2/N)
	// https://www.medcalc.org/manual/chi-square-table.php p=0.1
	var chiTable = map[int]float64{9: 14.68, 99: 117.407}

	t.Run("sort", func(t *testing.T) {
		var (
			i      uint64
			nodes  [size]uint64
			counts = make(map[Hashable]uint64, size)
			key    = make([]byte, 16)
		)

		for i = 0; i < size; i++ {
			nodes[i] = i
		}

		for i = 0; i < keys; i++ {
			binary.BigEndian.PutUint64(key, i+size)
			nodesHashed := wrapUint64(nodes[:])
			Sort(nodesHashed, WrapBytes(key))
			counts[nodesHashed[0]]++
		}

		var chi2 float64
		mean := float64(keys) / float64(size)
		delta := mean * percent
		for node, count := range counts {
			d := mean - float64(count)
			chi2 += math.Pow(float64(count)-mean, 2) / mean
			require.True(t, d < delta && (0-d) < delta,
				"Node %d received %d keys, expected %.0f (+/- %.2f)", node, count, mean, delta)
		}
		require.True(t, chi2 < chiTable[size-1],
			"Chi2 condition for .9 is not met (expected %.2f <= %.2f)", chi2, chiTable[size-1])
	})

	t.Run("hash collision", func(t *testing.T) {
		var (
			i      uint64
			counts = make(map[uint64]uint64)
			key    = make([]byte, 16)
		)

		for i = 0; i < keys; i++ {
			binary.BigEndian.PutUint64(key, i+size)
			hash := WrapBytes(key).Hash()
			counts[hash]++
		}

		for node, count := range counts {
			if count > 1 {
				t.Errorf("Node %d received %d keys", node, count)
			}
		}
	})
}
