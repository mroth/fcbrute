package main

import (
	"context"
	"math/rand"
	"strings"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

func Worker(ctx context.Context, r *rand.Rand, target string, results chan<- secp256k1.PrivateKey) {
	for {
		if ctx.Err() != nil {
			return
		}

		key, ok := workerIteration(r, target)
		if ok {
			results <- key
		}
	}
}

// workerIteration is a single iteration of the worker process, which uses r as
// a source of randomness to generate a new keypair, and check whether the
// string address matches the prefix target.
//
// It's isolated here as a means to easy holistic overall throughput
// measuring (we benchmark individual components in generate.go for
// performance tweaking).
func workerIteration(r *rand.Rand, target string) (secp256k1.PrivateKey, bool) {
	key, err := GenerateKeyInsecure(r)
	if err != nil {
		panic(err)
	}

	pubkey := key.PubKey()
	data := pubkey.SerializeUncompressed()
	addr, err := NewAddress(data)
	if err != nil {
		panic(err)
	}

	if strings.HasPrefix(addr, target) {
		return key, true
	}
	return key, false
}
