module github.com/utxorpc/go-sdk

go 1.21

// XXX: uncomment when testing local changes to spec, after generate
// replace github.com/utxorpc/go-codegen => ../go-codegen

require (
	connectrpc.com/connect v1.16.2
	github.com/utxorpc/go-codegen v0.9.0
	golang.org/x/net v0.29.0
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)
