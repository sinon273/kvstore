package storage

import (
	"strconv"
	"testing"
)

func BenchmarkSet(b *testing.B) {
	s := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Set("key", "value", 0)
	}
}

func BenchmarkGet(b *testing.B) {
	s := New()
	s.Set("key", "value", 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Get("key")
	}
}

func BenchmarkSetParallel(b *testing.B) {
	s := New()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Set("key", "value", 0)
		}
	})
}

func BenchmarkGetParallel(b *testing.B) {
	s := New()
	s.Set("key", "value", 0)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Get("key")
		}
	})
}

func BenchmarkGetParallelUniqueKeys(b *testing.B) {
	s := New()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "key" + strconv.Itoa(i)
			s.Set(key, "value", 0)
			s.Get(key)
			i++
		}
	})
}
