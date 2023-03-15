package main

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkWorkerIteration(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < b.N; i++ {
		workerIteration(r, "foobar")
	}
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "keys/sec")
}

func BenchmarkWorkerIteration_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			workerIteration(r, "foobar")
		}
	})
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "keys/sec")
}
