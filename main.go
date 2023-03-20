package main

import (
	"context"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"runtime"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

var (
	seed    = flag.Int64("seed", seedStart(), "random seed")
	workers = flag.Int("workers", runtime.NumCPU(), "number of workers")
)

// use the cryptographically secure RNG to initially seed the insecure RNGs
func seedStart() int64 {
	n, err := crand.Int(crand.Reader, big.NewInt(1<<63-1))
	if err != nil {
		panic(err)
	}
	return n.Int64()
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	target := flag.Arg(0)
	log.Printf("searching for %s with %d workers...", target, *workers)

	ctx := context.TODO()
	results := make(chan secp256k1.PrivateKey)

	for i := 0; i < *workers; i++ {
		r := rand.New(rand.NewSource(*seed + int64(i)))
		go Worker(ctx, r, target, results)
	}

	key := <-results

	pubkey := key.PubKey()
	pubkeydata := pubkey.SerializeUncompressed()
	smasher := NewAddressSmasher()
	smasher.Write(pubkeydata)
	addr := smasher.String()
	log.Printf("Found: %s\n", addr)

	keydata := key.Serialize()
	fmt.Println("private key (hex):", hex.EncodeToString(keydata))
	fmt.Println("private key (b64):", base64.StdEncoding.EncodeToString(keydata))
}

/*
For future use in importing into Lotus.  Need to figure out proper format.

// https://github.com/filecoin-project/lotus/blob/4fd81e0c58ba075c1410abf29469bc0ba1081124/chain/types/keystore.go

// KeyType defines a type of a key
type KeyType string

const (
	KTBLS             KeyType = "bls"
	KTSecp256k1       KeyType = "secp256k1"
	KTSecp256k1Ledger KeyType = "secp256k1-ledger"
	KTDelegated       KeyType = "delegated"
)

// KeyInfo is used for storing keys in KeyStore
// https://github.com/filecoin-project/venus/blob/73746cec80e43afd292f7d1ca68fb4774f783946/pkg/wallet/key/keyinfo.go#L24
type KeyInfo struct {
	PrivateKey []byte  `json:"privateKey"`
	Type       KeyType `json:"type"`
}
*/
