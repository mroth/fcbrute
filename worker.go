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
			results <- *key
		}
	}
}

