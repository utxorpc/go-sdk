module github.com/utxorpc/go-sdk

go 1.24.0

toolchain go1.24.1

// XXX: uncomment when testing local changes to spec, after generate
// replace github.com/utxorpc/go-codegen => ../go-codegen

require (
	connectrpc.com/connect v1.18.1
	github.com/blinklabs-io/gouroboros v0.135.2
	github.com/utxorpc/go-codegen v0.17.0
	golang.org/x/net v0.44.0
	google.golang.org/protobuf v1.36.9
)

require (
	github.com/bits-and-blooms/bitset v1.20.0 // indirect
	github.com/blinklabs-io/plutigo v0.0.12 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.5 // indirect
	github.com/btcsuite/btcd/btcutil v1.1.6 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.1.0 // indirect
	github.com/consensys/gnark-crypto v0.19.0 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.1.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
)
