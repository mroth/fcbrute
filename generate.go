package main

import (
	cryptorand "crypto/rand"
	mathrand "math/rand"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// PUBKEY GENERATION

// GenerateKey is the ported key generation from github.com/decred/dcrd/dcrec/secp256k1/v4.
// It is preserved here for reference convenience and benchmark comparison.
func GenerateKey() (*secp256k1.PrivateKey, error) {
	// The group order is close enough to 2^256 that there is only roughly a 1
	// in 2^128 chance of generating an invalid private key, so this loop will
	// virtually never run more than a single iteration in practice.
	var key secp256k1.PrivateKey
	var b32 [32]byte
	for valid := false; !valid; {
		if _, err := cryptorand.Read(b32[:]); err != nil {
			return nil, err
		}

		// The private key is only valid when it is in the range [1, N-1], where
		// N is the order of the curve.
		overflow := key.Key.SetBytes(&b32)
		valid = (key.Key.IsZeroBit() | overflow) == 0
	}
	zeroArray32(&b32)
	return &key, nil
}

// zeroArray32 zeroes the provided 32-byte buffer.
func zeroArray32(b *[32]byte) {
	copy(b[:], zero32[:])
}

var (
	// zero32 is an array of 32 bytes used for the purposes of zeroing and is
	// defined here to avoid extra allocations.
	zero32 = [32]byte{}
)

// GenerateKeyInsecure generates a secp256k1.PrivateKey, using insecure sources
// of randomness for greater speed and parallelization.
func GenerateKeyInsecure(r *mathrand.Rand) (secp256k1.PrivateKey, error) {
	// The group order is close enough to 2^256 that there is only roughly a 1
	// in 2^128 chance of generating an invalid private key, so this loop will
	// virtually never run more than a single iteration in practice.
	var key secp256k1.PrivateKey
	var b32 [32]byte
	for valid := false; !valid; {
		if _, err := r.Read(b32[:]); err != nil {
			return key, err
		}

		// The private key is only valid when it is in the range [1, N-1], where
		// N is the order of the curve.
		overflow := key.Key.SetBytes(&b32)
		valid = (key.Key.IsZeroBit() | overflow) == 0
	}
	// TODO: dcrec does zeroArray32(&b32) -- figure out why needs to be zeroed.
	// assuming this is just a memory enclave thing, can safely ignore.
	return key, nil
}
