package main

import (
	"bytes"
	"context"
	"math/rand"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

type worker struct {
	ctx     context.Context
	r       *rand.Rand
	smasher *AddressSmasher
	target  []byte

	bufPubkeyData [65]byte
}

func Worker(ctx context.Context, r *rand.Rand, target string, results chan<- secp256k1.PrivateKey) {
	w := worker{
		ctx:     ctx,
		r:       r,
		smasher: NewAddressSmasher(),
		target:  []byte(target),
	}

	for {
		if w.ctx.Err() != nil {
			return
		}

		key, ok := w.iterate()
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
func (w *worker) iterate() (secp256k1.PrivateKey, bool) {
	key, err := GenerateKeyInsecure(w.r)
	if err != nil {
		panic(err)
	}

	WritePubKeyBytes(&key, &w.bufPubkeyData)
	w.smasher.Write(w.bufPubkeyData[:])
	b := w.smasher.peekPayloadStringBytes()

	if bytes.HasPrefix(b, w.target) {
		return key, true
	}
	return key, false
}
