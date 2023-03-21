# fcbrute ⨎:hammer:

`fcbrute` will discover ⨎ Filecoin Protocol 1 (secp256k1) keypairs that result
in a public address string with a given prefix. This allows the creation of
"vanity" burner wallets.

The implementation does not rely on CGO, and performs zero memory heap
allocations during discovery while utilizing all CPU cores. A GPU-based
implementation would likely result in significantly faster performance. However,
on my laptop, the current testing throughput is ~0.31M keypairs/sec, which is
usually sufficient to find a 5-character prefix (my use case) in approximately
30 seconds, making it adequate for my current needs.

If you enjoy it, you can tip me FIL at `f1mroth3vxg3rqpetqx4pg2avaldxuwagi5dmlz4q`
(see what I did there?)

## Usage

    fcbrute -workers=8 abcd
    2023/03/19 17:45:06 searching for abcd with 8 workers...
    2023/03/19 17:45:08 Found: f1abcdblsbpha5ewj6ksm6bhhtl7dcky6xsvng7wq
    private key (hex): <redacted>
    private key (b64): <redacted>

## Status

Proof of concept.