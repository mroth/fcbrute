package main

import (
	"encoding/base32"
	"hash"

	"golang.org/x/crypto/blake2b"
)

// we want the address string, but avoiding any allocations.
// we also dont NEED the checksum for comparisons (I think, so we can probably avoid that...)
// technically we probably only need to encode the first N bytes that are necessary for our comparison even!

// Network represents which network an address belongs to.
type Network = byte

const (
	Mainnet Network = 'f' // Mainnet is the main network.
	Testnet Network = 't' // Testnet is the test network.
)

// Protocol represents which protocol an address uses.
type Protocol = byte

const (
	ID        Protocol = iota // ID represents the address ID protocol.
	SECP256K1                 // SECP256K1 represents the address SECP256K1 protocol.
	Actor                     // Actor represents the address Actor protocol.
	BLS                       // BLS represents the address BLS protocol.
	Delegated                 // Delegated represents the delegated (f4) address protocol.
	Unknown   = Protocol(255)
)

const encodeStd = "abcdefghijklmnopqrstuvwxyz234567"

// AddressSmasher is an optimized structure to generate lots of address strings
// with minimal wasted operations. It utilizes underlying buffers to minimize
// heap allocations and allow subsections of the address string to be calculated
// independently of eachother.
//
// It is intentionally not thread safe, utilize a local AddressSmasher for each
// thread.
//
// Currently AddressSmasher is hardcoded to only handle Protocol 1 addresses.
type AddressSmasher struct {
	/* variable settings */
	network  Network
	protocol Protocol

	/* local hashers and encoders */
	addressHasher  hash.Hash        // blake2b-160(data = raw uncompressed pubkey bytes) -> payload
	checksumHasher hash.Hash        // blake2b-32(data = protocol byte + payload) -> checksum
	low32          *base32.Encoding // lowercase base32 charset with no padding

	/* internal buffers */
	bufAddrHash []byte // payload, more or less
	bufChecksum []byte

	/*
		Internal byte representation of address string buffer.

			+-----------------+------------------------+------------------------+
			| prefix [2 byte] | b32(payload) [32 byte] | b32(checksum) [7 byte] |
			+-----------------+------------------------+------------------------+

		Note the base32 works out here to be equivalent such that:

			base32(payload + checksum) == base32(payload) + base32(checksum)

		This means we can calculate each independently to the same underlying
		buffer array at appropriate offsets depending on which we need.
	*/
	bufString []byte
}

const (
	prefixLen          = 2
	payloadLen         = 20 // hash160, 160 bits   = 20 bytes
	payloadLenEncoded  = 32 // base32raw(20 bytes) = 32 bytes
	checksumLen        = 4  // hash32, 32 bits     = 4 bytes
	checksumLenEncoded = 7  // base32raw(4 bytes)  = 7 bytes
)

func NewAddressSmasher() *AddressSmasher {

	hash160, err := blake2b.New(20, nil)
	if err != nil {
		panic(err)
	}

	hash32, err := blake2b.New(4, nil)
	if err != nil {
		panic(err)
	}

	return &AddressSmasher{
		network:  Mainnet,
		protocol: SECP256K1,

		addressHasher:  hash160,
		checksumHasher: hash32,
		low32:          base32.NewEncoding(encodeStd).WithPadding(base32.NoPadding),
		bufAddrHash:    make([]byte, 0, payloadLen),
		bufChecksum:    make([]byte, 0, checksumLen),
		bufString:      make([]byte, prefixLen+payloadLenEncoded+checksumLenEncoded),
	}
}

// Write inputs the public key, producing the payload.  Subsequent writes will
// overwrite the previous payload.
func (a *AddressSmasher) Write(pubkey []byte) {
	a.addressHasher.Reset()
	_, _ = a.addressHasher.Write(pubkey)

	a.bufAddrHash = a.bufAddrHash[:0] // reset
	a.bufAddrHash = a.addressHasher.Sum(a.bufAddrHash)
}

func (a *AddressSmasher) String() string {
	a.calcStrRep_Prefix()
	a.calcStrRep_Payload()
	a.calcStrRep_Checksum()
	return string(a.bufString)
}

// calculate and peek at the bytes of the payload portion of the address string
func (a *AddressSmasher) peekPayloadStringBytes() []byte {
	a.calcStrRep_Payload()
	return a.bufString[2:34]
}

func (a *AddressSmasher) calcStrRep_Prefix() {
	dst := a.bufString[:2]
	dst[0] = a.network
	dst[1] = a.protocol + 48 // ASCII offset for 0 digit
}

func (a *AddressSmasher) calcStrRep_Payload() {
	dst := a.bufString[2:34]
	a.low32.Encode(dst, a.bufAddrHash)
}

func (a *AddressSmasher) calcStrRep_Checksum() {
	a.checksumHasher.Reset()
	a.checksumHasher.Write([]byte{a.protocol})
	a.checksumHasher.Write(a.bufAddrHash)
	a.bufChecksum = a.bufChecksum[:0]
	a.bufChecksum = a.checksumHasher.Sum(a.bufChecksum)

	dst := a.bufString[34:]
	a.low32.Encode(dst, a.bufChecksum)
}
