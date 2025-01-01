package csmap

import (
	"sync"
	"testing"
)

func TestCSMap_SetAndGet(t *testing.T) {
	csMap := NewCSMap[string, int](10)

	csMap.Set("key1", 1)
	csMap.Set("key2", 2)

	value, ok := csMap.Get("key1")
	if !ok || value != 1 {
		t.Errorf("Expected value 1 for key 'key1', got %v, ok: %v", value, ok)
	}

	value, ok = csMap.Get("key2")
	if !ok || value != 2 {
		t.Errorf("Expected value 2 for key 'key2', got %v, ok: %v", value, ok)
	}

	value, ok = csMap.Get("key3")
	if ok {
		t.Errorf("Expected key 'key3' to not exist, but it does with value %v", value)
	}
}

func TestCSMap_Delete(t *testing.T) {
	csMap := NewCSMap[string, int](10)

	csMap.Set("key1", 1)
	csMap.Delete("key1")

	_, ok := csMap.Get("key1")
	if ok {
		t.Errorf("Expected key 'key1' to be deleted, but it still exists")
	}
}

func TestCSMap_ConcurrentAccess(t *testing.T) {
	csMap := NewCSMap[int, string](10)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			csMap.Set(i, "value")
		}(i)
	}

	wg.Wait()

	for i := 0; i < 100; i++ {
		value, ok := csMap.Get(i)
		if !ok || value != "value" {
			t.Errorf("Expected value 'value' for key %d, got %v, ok: %v", i, value, ok)
		}
	}
}

func TestCSMap_HashCollision(t *testing.T) {
	csMap := NewCSMap[string, string](2)

	csMap.Set("key1", "value1")
	csMap.Set("key2", "value2")

	value, ok := csMap.Get("key1")
	if !ok || value != "value1" {
		t.Errorf("Expected 'value1' for 'key1', got %v, ok: %v", value, ok)
	}

	value, ok = csMap.Get("key2")
	if !ok || value != "value2" {
		t.Errorf("Expected 'value2' for 'key2', got %v, ok: %v", value, ok)
	}
}

// cpu: Intel(R) Core(TM) i7-7700 CPU @ 3.60GHz
// BenchmarkCSMapSet
// BenchmarkCSMapSet-8      6084016               176.5 ns/op            59 B/op            0 allocs/op
func BenchmarkCSMapSet(b *testing.B) {
	csm := NewCSMap[int, int](32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		csm.Set(i, i)
	}
	b.ReportAllocs()
}

// cpu: Intel(R) Core(TM) i7-7700 CPU @ 3.60GHz
// BenchmarkCSMapGet
// BenchmarkCSMapGet-8     34310907                32.02 ns/op            0 B/op            0 allocs/op
func BenchmarkCSMapGet(b *testing.B) {
	csm := NewCSMap[int, int](32)
	for i := 0; i < 1000; i++ {
		csm.Set(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		csm.Get(i % 1000)
	}
	b.ReportAllocs()
}

// cpu: Intel(R) Core(TM) i7-7700 CPU @ 3.60GHz
// BenchmarkNativeMapSet
// BenchmarkNativeMapSet-8          8071980               167.4 ns/op            87 B/op          0 allocs/op
func BenchmarkNativeMapSet(b *testing.B) {
	nativeMap := make(map[int]int)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nativeMap[i] = i
	}
	b.ReportAllocs()
}

// cpu: Intel(R) Core(TM) i7-7700 CPU @ 3.60GHz
// BenchmarkNativeMapGet
// BenchmarkNativeMapGet-8         80680946                16.70 ns/op            0 B/op          0 allocs/op
func BenchmarkNativeMapGet(b *testing.B) {
	nativeMap := make(map[int]int)
	for i := 0; i < 1000; i++ {
		nativeMap[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = nativeMap[i%1000]
	}
	b.ReportAllocs()
}
