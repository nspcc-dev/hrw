# Golang HRW implementation

[![Build Status](https://travis-ci.org/nspcc-dev/hrw.svg?branch=master)](https://travis-ci.org/nspcc-dev/hrw)
[![codecov](https://codecov.io/gh/nspcc-dev/hrw/badge.svg)](https://codecov.io/gh/nspcc-dev/hrw)
[![Report](https://goreportcard.com/badge/github.com/nspcc-dev/hrw)](https://goreportcard.com/report/github.com/nspcc-dev/hrw)
[![GitHub release](https://img.shields.io/github/release/nspcc-dev/hrw.svg)](https://github.com/nspcc-dev/hrw)

[Rendezvous or highest random weight](https://en.wikipedia.org/wiki/Rendezvous_hashing) (HRW) hashing is an algorithm that allows clients to achieve distributed agreement on a set of k options out of a possible set of n options. A typical application is when clients need to agree on which sites (or proxies) objects are assigned to. When k is 1, it subsumes the goals of consistent hashing, using an entirely different method.

## Install

`go get github.com/nspcc-dev/hrw`

## Benchmark:

```
BenchmarkSort_fnv_10-8                           5000000               354 ns/op             224 B/op          3 allocs/op
BenchmarkSort_fnv_100-8                           300000              5103 ns/op            1856 B/op          3 allocs/op
BenchmarkSort_fnv_1000-8                           10000            115874 ns/op           16448 B/op          3 allocs/op
BenchmarkSortByIndex_fnv_10-8                    3000000               562 ns/op             384 B/op          7 allocs/op
BenchmarkSortByIndex_fnv_100-8                    200000              5819 ns/op            2928 B/op          7 allocs/op
BenchmarkSortByIndex_fnv_1000-8                    10000            125859 ns/op           25728 B/op          7 allocs/op
BenchmarkSortByValue_fnv_10-8                    2000000              1056 ns/op             544 B/op         17 allocs/op
BenchmarkSortByValue_fnv_100-8                    200000              9593 ns/op            4528 B/op        107 allocs/op
BenchmarkSortByValue_fnv_1000-8                    10000            109272 ns/op           41728 B/op       1007 allocs/op

BenchmarkSortByWeight_fnv_10-8                   3000000               500 ns/op             320 B/op          4 allocs/op
BenchmarkSortByWeight_fnv_100-8                   200000              8257 ns/op            2768 B/op          4 allocs/op
BenchmarkSortByWeight_fnv_1000-8                   10000            197938 ns/op           24656 B/op          4 allocs/op
BenchmarkSortByWeightIndex_fnv_10-8              2000000               760 ns/op             480 B/op          8 allocs/op
BenchmarkSortByWeightIndex_fnv_100-8              200000              9191 ns/op            3840 B/op          8 allocs/op
BenchmarkSortByWeightIndex_fnv_1000-8              10000            208204 ns/op           33936 B/op          8 allocs/op
BenchmarkSortByWeightValue_fnv_10-8              1000000              1095 ns/op             640 B/op         18 allocs/op
BenchmarkSortByWeightValue_fnv_100-8              200000             12291 ns/op            5440 B/op        108 allocs/op
BenchmarkSortByWeightValue_fnv_1000-8              10000            145125 ns/op           49936 B/op       1008 allocs/op

```

## Example

```go
package main

import (
	"fmt"
	
	"github.com/nspcc-dev/hrw"
)

func main() {
	// given a set of servers
	servers := []string{
		"one.example.com",
		"two.example.com",
		"three.example.com",
		"four.example.com",
		"five.example.com",
		"six.example.com",
	}

	// HRW can consistently select a uniformly-distributed set of servers for
	// any given key
	var (
		key = []byte("/examples/object-key")
		h   = hrw.Hash(key)
	)

	hrw.SortSliceByValue(servers, h)
	for id := range servers {
		fmt.Printf("trying GET %s%s\n", servers[id], key)
	}

	// Output:
	// trying GET four.example.com/examples/object-key
	// trying GET three.example.com/examples/object-key
	// trying GET one.example.com/examples/object-key
	// trying GET two.example.com/examples/object-key
	// trying GET six.example.com/examples/object-key
	// trying GET five.example.com/examples/object-key
}
```