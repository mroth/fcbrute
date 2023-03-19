package main

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkWorkerIteration(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	w := worker{
		ctx:     context.Background(),
		r:       r,
		smasher: NewAddressSmasher(),
		target:  []byte("foobar"),
	}

	for i := 0; i < b.N; i++ {
		w.iterate()
	}
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "keys/sec")
}

func BenchmarkWorkerIteration_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		w := worker{
			ctx:     context.Background(),
			r:       r,
			smasher: NewAddressSmasher(),
			target:  []byte("foobar"),
		}

		for pb.Next() {
			w.iterate()
		}
	})
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "keys/sec")
}
