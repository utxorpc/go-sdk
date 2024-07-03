package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/cardano"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/query"
	utxorpc "github.com/utxorpc/go-sdk"
)

func main() {

	ctx := context.Background()
	baseUrl := "https://preview.utxorpc-v0.demeter.run"
	// set API key for demeter
	client := utxorpc.CreateUtxoRPCClient(baseUrl,
		// set API key for demeter
		utxorpc.WithHeaders(map[string]string{
			"dmtr-api-key": "dmtr_apikey...",
		}),
	)

	// Set mode to "readParams", "readUtxos", "searchUtxos" to select the desired example.
	var mode string = "readParams"

	switch mode {
	case "readParams":
		readParams(ctx, client)
	case "readUtxos":
		readUtxos(ctx, client)
	case "searchUtxos":
		searchUtxos(ctx, client)
	default:
		fmt.Println("Unknown mode:", mode)
	}
}

func readParams(ctx context.Context, client *utxorpc.UtxorpcClient) {
	req := connect.NewRequest(&query.ReadParamsRequest{})
	client.AddHeadersToRequest(req)

	fmt.Println("Connecting to utxorpc host:", client.URL())
	resp, err := client.Query.ReadParams(ctx, req)
	if err != nil {
		utxorpc.HandleError(err)
	}
	fmt.Printf("Response: %+v\n", resp)
}

func readUtxos(ctx context.Context, client *utxorpc.UtxorpcClient) {
	txHash, err := hex.DecodeString("3394533cb02fb71b062690d85bbe9d79a7b6f8f4c1b92b0e728fe7b93a1440c9")
	if err != nil {
		log.Fatalf("failed to decode hex string: %v", err)
	}
	txoRef := &query.TxoRef{
		Hash: txHash,
	}

	req := connect.NewRequest(&query.ReadUtxosRequest{
		Keys: []*query.TxoRef{txoRef},
	})
	client.AddHeadersToRequest(req)
	fmt.Println("connecting to utxorpc host:", client.URL())
	resp, err := client.Query.ReadUtxos(ctx, req)
	if err != nil {
		utxorpc.HandleError(err)
	}
	fmt.Printf("Response: %+v\n", resp)
}

func searchUtxos(ctx context.Context, client *utxorpc.UtxorpcClient) {

	rawAddress := "00fc26cca7cc4b67032c38ed6a568237cc05ba73e31d375198ffe0e7f56c4506809de7d19eb9b619908e2027471e97fe953b3234a476258560"
	exactAddress, err := hex.DecodeString(rawAddress)
	if err != nil {
		log.Fatalf("failed to decode hex string address: %v", err)
	}

	req := connect.NewRequest(&query.SearchUtxosRequest{
		Predicate: &query.UtxoPredicate{
			Match: &query.AnyUtxoPattern{
				UtxoPattern: &query.AnyUtxoPattern_Cardano{
					Cardano: &cardano.TxOutputPattern{
						Address: &cardano.AddressPattern{
							ExactAddress: exactAddress,
						},
						Asset: &cardano.AssetPattern{
							// Populate the fields as necessary
						},
					},
				},
			},
		},
	})
	client.AddHeadersToRequest(req)

	fmt.Println("connecting to utxorpc host:", client.URL())
	resp, err := client.Query.SearchUtxos(ctx, req)
	if err != nil {
		utxorpc.HandleError(err)
	}
	fmt.Printf("Response: %+v\n", resp)
}
