module github.com/utxorpc/go-sdk

go 1.22

toolchain go1.22.8

// XXX: uncomment when testing local changes to spec, after generate
// replace github.com/utxorpc/go-codegen => ../go-codegen

require (
	connectrpc.com/connect v1.17.0
	github.com/blinklabs-io/gouroboros v0.105.0
	github.com/utxorpc/go-codegen v0.12.0
	golang.org/x/net v0.31.0
	google.golang.org/protobuf v1.35.2
)

require (
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.29.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/text v0.20.0 // indirect
)
