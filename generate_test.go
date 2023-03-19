package main

import (
	"bytes"
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

/*
We now use WritePubKeyBytes to avoid these two allocating secp256k1 functions,
but preserve our benchmarks below which allowed us to see their allocations.

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
*/

// make sure we get the same resulting bytes as using the secp256k1 module directly,
// by testing a set of random keys.
func TestWritePubKeyBytes(t *testing.T) {
	const samples = 10
	for i := 0; i < samples; i++ {
		key, err := GenerateKey()
		if err != nil {
			t.Fatal(err)
		}

		// our method
		var data [65]byte
		WritePubKeyBytes(key, &data)
		got := data[:]

		// secp256k1 library with allocations
		want := key.PubKey().SerializeUncompressed()

		// compare
		if !bytes.Equal(got, want) {
			t.Errorf("got %+x want %+x", got, want)
		}
	}
}

func BenchmarkWritePubKeyBytes(b *testing.B) {
	key, err := GenerateKey()
	if err != nil {
		b.Fatal(err)
	}

	var data [65]byte
	for i := 0; i < b.N; i++ {
		WritePubKeyBytes(key, &data)
	}
}
