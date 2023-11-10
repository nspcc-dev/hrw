package hrw

import "testing"

func BenchmarkSort_fnv_10(b *testing.B) {
	_ = benchmarkSort(b, 10, testKey)
}

func BenchmarkSort_fnv_100(b *testing.B) {
	_ = benchmarkSort(b, 100, testKey)
}

func BenchmarkSort_fnv_1000(b *testing.B) {
	_ = benchmarkSort(b, 1000, testKey)
}

func BenchmarkSortByWeight_fnv_10(b *testing.B) {
	_ = benchmarkSortByWeight(b, 10, testKey)
}

func BenchmarkSortByWeight_fnv_100(b *testing.B) {
	_ = benchmarkSortByWeight(b, 100, testKey)
}

func BenchmarkSortByWeight_fnv_1000(b *testing.B) {
	_ = benchmarkSortByWeight(b, 1000, testKey)
}

func benchmarkSort(b *testing.B, n int, object []byte) uint64 {
	servers := make([]hashableUint64, n)
	for i := uint64(0); i < uint64(len(servers)); i++ {
		servers[i] = hashableUint64(i)
	}

	oHash := hashableUint64(WrapBytes(object).Hash())

	b.ResetTimer()
	b.ReportAllocs()

	var x uint64
	for i := 0; i < b.N; i++ {
		Sort(servers, oHash)
		x += servers[0].Hash()
	}
	return x
}

func benchmarkSortByWeight(b *testing.B, n int, object []byte) uint64 {
	servers := make([]hashableUint64, n)
	weights := make([]float64, n)
	for i := uint64(0); i < uint64(len(servers)); i++ {
		weights[i] = float64(uint64(n)-i) / float64(n)
		servers[i] = hashableUint64(i)
	}

	oHash := hashableUint64(WrapBytes(object).Hash())

	b.ResetTimer()
	b.ReportAllocs()

	var x uint64
	for i := 0; i < b.N; i++ {
		SortWeighted(servers, weights, oHash)
		x += servers[0].Hash()
	}
	return x
}
