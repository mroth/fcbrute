package main

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkGenerateKey_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GenerateKey()
		}
	})
}

func BenchmarkGenerateKeyInsecure_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			GenerateKeyInsecure(r)
		}
	})
}

func BenchmarkGeneratePubkey(b *testing.B) {
	key, err := GenerateKey()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		key.PubKey()
	}
}

func BenchmarkPubkeySerialize(b *testing.B) {
	key, err := GenerateKey()
	if err != nil {
		b.Fatal(err)
	}
	pubkey := key.PubKey()

	for i := 0; i < b.N; i++ {
		pubkey.SerializeUncompressed()
	}
}

func BenchmarkNewAddress(b *testing.B) {
	key, err := GenerateKey()
	if err != nil {
		b.Fatal(err)
	}
	pubkey := key.PubKey()
	pubkeydata := pubkey.SerializeUncompressed()

	for i := 0; i < b.N; i++ {
		_, _ = NewAddress(pubkeydata)
	}
}
