package main

import (
	"encoding/base32"
	"encoding/base64"
	"testing"
)

func TestBase32(t *testing.T) {
	var low32 = base32.NewEncoding(encodeStd).WithPadding(base32.NoPadding)

	const (
		s1 = "abcdefghijklmnopqrst"
		s2 = "1234"
	)

	t.Logf("b32(s1) \t%v", low32.EncodeToString([]byte(s1)))
	t.Logf("b32(s2) \t%v", low32.EncodeToString([]byte(s2)))
	t.Logf("b32(s1+s2)\t%v", low32.EncodeToString([]byte(s1+s2)))
	t.Logf("b32(s1)+b32(s2)\t%v%v", low32.EncodeToString([]byte(s1)), low32.EncodeToString([]byte(s2)))
}

func TestAddressSmasher_String(t *testing.T) {
	// some randomly generated test cases pushed through filecoin-project/go-address to test we get same results
	var testcases = []struct {
		data string // base64-encoded 65 byte uncompressed public key
		want string // expected address string
	}{
		{
			data: "BPNvb8F8oBntQkruKWHxKe/K5l77C8dcIOvgKqh7iKaCDUCZJIfoU59bQn+hOt5Ymh4kjMbjiBbHsd2Sg1DScQQ",
			want: "f1pkqjza7n7ds7zw2tvytehzlzj57avoqtgmuppja",
		},
		{
			data: "BFzcGbvUmcACHAFdSggUIVe6QphGe7AZLLgavr8+EtZZenjwbO/hQfYkH5/rzs5Y3zWmyiFHpRVhapqazRMQOao",
			want: "f1ez5lqpu326eggq5c2ipxcilzaa2zq6rncb4xsfa",
		},
		{
			data: "BP0ZkvMAYfvSx2uFjVyr/wj4Zk5If38WI3wEdbvVmB84VoZ2fGNr6rpc5gWyDHgqHPcfCS8pJnEr8o/4jl4pnX8",
			want: "f1asisadfq6ifufxj7bfhnrxlvoy6ruvnwi46dbgy",
		},
		{
			data: "BOvXSllhID4gGLdCYX8gkBSj8nIBJgL0eRCWHbC08L0omtBDK7ejKPLRinACIu8uLEc7YpAOHBqNjRN+PZOZQDE",
			want: "f1y7zonhkf3tqoy6elxm3ltdmh75xnfcy3okbsmti",
		},
		{
			data: "BP4IMzQWy550KuFuNU1rxNNJWEtbBsSnrEzhMounJKPc33EjQT4W5nunuSPMTsjhBHi4hFR9b5nK8xl0/WzrAlY",
			want: "f1u27y32tem2pvlaqvnuguvtweoom7v2ilwx34eeq",
		},
	}

	smasher := NewAddressSmasher()
	for _, tc := range testcases {
		keydata, err := base64.RawStdEncoding.DecodeString(tc.data)
		if err != nil {
			t.Fatal(err)
		}
		smasher.Write(keydata)
		got := smasher.String()
		if tc.want != got {
			t.Errorf("want %v got %v", tc.want, got)
		}
	}
}

func BenchmarkAddressSmasher_Write(b *testing.B) {
	b.ReportAllocs()

	smasher := NewAddressSmasher()
	pubkeydata := generatePubKeyData(b)

	for i := 0; i < b.N; i++ {
		smasher.Write(pubkeydata)
	}
}

func BenchmarkAddressSmasher_String(b *testing.B) {
	b.ReportAllocs()

	smasher := NewAddressSmasher()
	pubkeydata := generatePubKeyData(b)
	smasher.Write(pubkeydata)

	for i := 0; i < b.N; i++ {
		_ = smasher.String()
	}
}

func BenchmarkAddressSmasher_peekPayloadStringBytes(b *testing.B) {
	b.ReportAllocs()

	smasher := NewAddressSmasher()
	pubkeydata := generatePubKeyData(b)
	smasher.Write(pubkeydata)

	for i := 0; i < b.N; i++ {
		_ = smasher.peekPayloadStringBytes()
	}
}

func generatePubKeyData(tb testing.TB) []byte {
	tb.Helper()

	key, err := GenerateKey()
	if err != nil {
		tb.Fatal(err)
	}
	pubkey := key.PubKey()
	pubkeydata := pubkey.SerializeUncompressed()

	return pubkeydata
}
