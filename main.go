package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/filecoin-project/go-address"
)

var (
	// TODO use unpredictable seed start
	seed    = flag.Int64("seed", time.Now().UnixNano(), "random seed")
	workers = flag.Int("workers", runtime.NumCPU(), "number of workers")
)

func main() {
	flag.Parse()
	target := "t1abc"

	ctx := context.TODO()
	results := make(chan secp256k1.PrivateKey)

	for i := 0; i < *workers; i++ {
		r := rand.New(rand.NewSource(*seed + int64(i)))
		go Worker(ctx, r, target, results)
	}

	key := <-results

	pubkey := key.PubKey()
	address.CurrentNetwork = address.Mainnet
	addr, err := NewAddress(pubkey.SerializeUncompressed())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found: %s\n", addr)

	keydata := key.Serialize()
	fmt.Println("hex private key:", hex.EncodeToString(keydata))
	fmt.Println("b64 private key:", base64.StdEncoding.EncodeToString(keydata))

	info := KeyInfo{
		Type:       KTSecp256k1,
		PrivateKey: keydata,
	}
	exportKey, err := json.Marshal(info)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("hex-lotus: %s\n", exportKey)
}

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
