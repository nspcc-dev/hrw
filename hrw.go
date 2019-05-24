// Package hrw implements Rendezvous hashing.
// http://en.wikipedia.org/wiki/Rendezvous_hashing.
package hrw

import (
	"encoding/binary"
	"reflect"
	"sort"

	"github.com/spaolacci/murmur3"
)

type (
	swapper func(i, j int)

	// Hasher interface used by SortSliceByValue
	Hasher interface{ Hash() uint64 }

	hashed struct {
		length int
		sorted []uint64
		weight []uint64
	}

	weighted struct {
		h      hashed
		normal []float64 // normalized input weights
	}
)

func weight(x uint64, y uint64) uint64 {
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

func (h hashed) Len() int           { return h.length }
func (h hashed) Less(i, j int) bool { return h.weight[i] < h.weight[j] }
func (h hashed) Swap(i, j int) {
	h.sorted[i], h.sorted[j] = h.sorted[j], h.sorted[i]
	h.weight[i], h.weight[j] = h.weight[j], h.weight[i]
}

func (w weighted) Len() int { return w.h.length }
func (w weighted) Less(i, j int) bool {
	// `maxUint64 - weight` makes least weight most valuable
	// it is necessary for operation with normalized values
	wi := float64(^uint64(0)-w.h.weight[i]) * w.normal[i]
	wj := float64(^uint64(0)-w.h.weight[j]) * w.normal[j]
	return wi > wj // higher weight must be placed lower to be first
}
func (w weighted) Swap(i, j int) { w.normal[i], w.normal[j] = w.normal[j], w.normal[i]; w.h.Swap(i, j) }

// Hash uses murmur3 hash to return uint64
func Hash(key []byte) uint64 {
	return murmur3.Sum64(key)
}

// Sort receive nodes and hash, and sort it by weight
func Sort(nodes []uint64, hash uint64) []uint64 {
	var (
		l = len(nodes)
		h = hashed{
			length: l,
			sorted: make([]uint64, 0, l),
			weight: make([]uint64, 0, l),
		}
	)

	for i, node := range nodes {
		h.sorted = append(h.sorted, uint64(i))
		h.weight = append(h.weight, weight(node, hash))
	}

	sort.Sort(h)
	return h.sorted
}

// SortByWeight receive nodes and hash, and sort it by weight
func SortByWeight(nodes []uint64, weights []uint64, hash uint64) []uint64 {
	var (
		maxWeight uint64

		l = len(nodes)
		w = weighted{
			h: hashed{
				length: l,
				sorted: make([]uint64, 0, l),
				weight: make([]uint64, 0, l),
			},
			normal: make([]float64, 0, l),
		}
	)

	// finding max weight to perform normalization
	for i := range weights {
		if maxWeight < weights[i] {
			maxWeight = weights[i]
		}
	}

	// if all nodes have 0-weights or weights are incorrect then sort uniformly
	if maxWeight == 0 || l != len(nodes) {
		return Sort(nodes, hash)
	}

	fMaxWeight := float64(maxWeight)
	for i, node := range nodes {
		w.h.sorted = append(w.h.sorted, uint64(i))
		w.h.weight = append(w.h.weight, weight(node, hash))
		w.normal = append(w.normal, float64(weights[i])/fMaxWeight)
	}
	sort.Sort(w)
	return w.h.sorted
}

// SortSliceByValue received []T and hash to sort by value-weight
func SortSliceByValue(slice interface{}, hash uint64) {
	rule := prepareRule(slice)
	if rule != nil {
		swap := reflect.Swapper(slice)
		rule = Sort(rule, hash)
		sortByRuleInverse(swap, uint64(len(rule)), rule)
	}
}

// SortSliceByWeightValue received []T, weights and hash to sort by value-weight
func SortSliceByWeightValue(slice interface{}, weight []uint64, hash uint64) {
	rule := prepareRule(slice)
	if rule != nil {
		swap := reflect.Swapper(slice)
		rule = SortByWeight(rule, weight, hash)
		sortByRuleInverse(swap, uint64(len(rule)), rule)
	}
}

// SortSliceByIndex received []T and hash to sort by index-weight
func SortSliceByIndex(slice interface{}, hash uint64) {
	length := uint64(reflect.ValueOf(slice).Len())
	swap := reflect.Swapper(slice)
	rule := make([]uint64, 0, length)
	for i := uint64(0); i < length; i++ {
		rule = append(rule, i)
	}
	rule = Sort(rule, hash)
	sortByRuleInverse(swap, length, rule)
}

// SortSliceByWeightIndex received []T, weights and hash to sort by index-weight
func SortSliceByWeightIndex(slice interface{}, weight []uint64, hash uint64) {
	length := uint64(reflect.ValueOf(slice).Len())
	swap := reflect.Swapper(slice)
	rule := make([]uint64, 0, length)
	for i := uint64(0); i < length; i++ {
		rule = append(rule, i)
	}
	rule = SortByWeight(rule, weight, hash)
	sortByRuleInverse(swap, length, rule)
}

func sortByRuleDirect(swap swapper, length uint64, rule []uint64) {
	done := make([]bool, length)
	for i := uint64(0); i < length; i++ {
		if done[i] {
			continue
		}
		for j := rule[i]; !done[rule[j]]; j = rule[j] {
			swap(int(i), int(j))
			done[j] = true
		}
	}
}

func sortByRuleInverse(swap swapper, length uint64, rule []uint64) {
	done := make([]bool, length)
	for i := uint64(0); i < length; i++ {
		if done[i] {
			continue
		}

		for j := i; !done[rule[j]]; j = rule[j] {
			swap(int(j), int(rule[j]))
			done[j] = true
		}
	}
}

func prepareRule(slice interface{}) []uint64 {
	t := reflect.TypeOf(slice)
	if t.Kind() != reflect.Slice {
		return nil
	}

	var (
		val    = reflect.ValueOf(slice)
		length = val.Len()
		rule   = make([]uint64, 0, length)
	)

	if length == 0 {
		return nil
	}

	switch slice := slice.(type) {
	case []int:
		var key = make([]byte, 16)
		for i := 0; i < length; i++ {
			binary.BigEndian.PutUint64(key, uint64(slice[i]))
			rule = append(rule, Hash(key))
		}
	case []uint:
		var key = make([]byte, 16)
		for i := 0; i < length; i++ {
			binary.BigEndian.PutUint64(key, uint64(slice[i]))
			rule = append(rule, Hash(key))
		}
	case []int8:
		for i := 0; i < length; i++ {
			key := byte(slice[i])
			rule = append(rule, Hash([]byte{key}))
		}
	case []uint8:
		for i := 0; i < length; i++ {
			key := slice[i]
			rule = append(rule, Hash([]byte{key}))
		}
	case []int16:
		var key = make([]byte, 8)
		for i := 0; i < length; i++ {
			binary.BigEndian.PutUint16(key, uint16(slice[i]))
			rule = append(rule, Hash(key))
		}
	case []uint16:
		var key = make([]byte, 8)
		for i := 0; i < length; i++ {
			binary.BigEndian.PutUint16(key, slice[i])
			rule = append(rule, Hash(key))
		}
	case []int32:
		var key = make([]byte, 16)
		for i := 0; i < length; i++ {
			binary.BigEndian.PutUint32(key, uint32(slice[i]))
			rule = append(rule, Hash(key))
		}
	case []uint32:
		var key = make([]byte, 16)
		for i := 0; i < length; i++ {
			binary.BigEndian.PutUint32(key, slice[i])
			rule = append(rule, Hash(key))
		}
	case []int64:
		var key = make([]byte, 32)
		for i := 0; i < length; i++ {
			binary.BigEndian.PutUint64(key, uint64(slice[i]))
			rule = append(rule, Hash(key))
		}
	case []uint64:
		var key = make([]byte, 32)
		for i := 0; i < length; i++ {
			binary.BigEndian.PutUint64(key, slice[i])
			rule = append(rule, Hash(key))
		}
	case []string:
		for i := 0; i < length; i++ {
			rule = append(rule, Hash([]byte(slice[i])))
		}

	default:
		if _, ok := val.Index(0).Interface().(Hasher); !ok {
			return nil
		}

		for i := 0; i < length; i++ {
			h := val.Index(i).Interface().(Hasher)
			rule = append(rule, h.Hash())
		}
	}
	return rule
}
