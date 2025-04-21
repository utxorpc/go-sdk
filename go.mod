module github.com/utxorpc/go-sdk

go 1.23.6

toolchain go1.24.1

// XXX: uncomment when testing local changes to spec, after generate
// replace github.com/utxorpc/go-codegen => ../go-codegen

require (
	connectrpc.com/connect v1.18.1
	github.com/blinklabs-io/gouroboros v0.117.0
	github.com/utxorpc/go-codegen v0.16.0
	golang.org/x/net v0.39.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/btcsuite/btcd/btcutil v1.1.6 // indirect
	github.com/fxamacker/cbor/v2 v2.8.0 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)
