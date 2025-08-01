module github.com/utxorpc/go-sdk

go 1.24.0

toolchain go1.24.1

// XXX: uncomment when testing local changes to spec, after generate
// replace github.com/utxorpc/go-codegen => ../go-codegen

require (
	connectrpc.com/connect v1.18.1
	github.com/blinklabs-io/gouroboros v0.129.0
	github.com/utxorpc/go-codegen v0.17.0
	golang.org/x/net v0.42.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/blinklabs-io/plutigo v0.0.2-0.20250717183329-b331a97fb319 // indirect
	github.com/btcsuite/btcd/btcutil v1.1.6 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
)
