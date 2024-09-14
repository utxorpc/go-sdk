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
	var mode string = "searchUtxos"

	switch mode {
	case "readParams":
		readParams(ctx, client)
	case "readUtxos":
		readUtxos(ctx, client, "71a7498f086d378ec5e558581286629b678be1dd65d5d4e2a5d634ba6fdf8299", 0)
	case "searchUtxos":
		searchUtxos(ctx, client, "60c0359ebb7d0688d79064bd118c99c8b87b5853e3af59245bb97e84d2")
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

	if resp.Msg.LedgerTip != nil {
		fmt.Printf("Ledger Tip: Slot: %d, Hash: %x\n", resp.Msg.LedgerTip.Slot, resp.Msg.LedgerTip.Hash)
	}
	if resp.Msg.Values != nil {
		fmt.Printf("Cardano: %+v\n", resp.Msg.Values)
	}
}

func readUtxos(ctx context.Context, client *utxorpc.UtxorpcClient, txHashStr string, txIndex uint32) {
	txHash, err := hex.DecodeString(txHashStr)
	if err != nil {
		log.Fatalf("failed to decode hex string: %v", err)
	}
	txoRef := &query.TxoRef{
		Hash:  txHash,
		Index: txIndex,
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

	if resp.Msg.LedgerTip != nil {
		fmt.Printf("Ledger Tip:\n  Slot: %d\n  Hash: %x\n", resp.Msg.LedgerTip.Slot, resp.Msg.LedgerTip.Hash)
	}

	for _, item := range resp.Msg.Items {
		fmt.Println("UTxO Data:")
		fmt.Printf("  Tx Hash: %x\n", item.TxoRef.Hash)
		fmt.Printf("  Output Index: %d\n", item.TxoRef.Index)
		fmt.Printf("  Native Bytes: %x\n", item.NativeBytes)
		if cardano := item.GetCardano(); cardano != nil {
			fmt.Println("  Cardano UTxO:")
			fmt.Printf("    Address: %s\n", cardano.Address)
			fmt.Printf("    Coin: %d\n", cardano.Coin)
		}
	}
}

func searchUtxos(ctx context.Context, client *utxorpc.UtxorpcClient, rawAddress string) {
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

	if resp.Msg.LedgerTip != nil {
		fmt.Printf("Ledger Tip:\n  Slot: %d\n  Hash: %x\n", resp.Msg.LedgerTip.Slot, resp.Msg.LedgerTip.Hash)
	}

	for _, item := range resp.Msg.Items {
		fmt.Println("UTxO Data:")
		fmt.Printf("  Tx Hash: %x\n", item.TxoRef.Hash)
		fmt.Printf("  Output Index: %d\n", item.TxoRef.Index)
		fmt.Printf("  Native Bytes: %x\n", item.NativeBytes)
		if cardano := item.GetCardano(); cardano != nil {
			fmt.Println("  Cardano UTxO:")
			fmt.Printf("    Address: %s\n", cardano.Address)
			fmt.Printf("    Coin: %d\n", cardano.Coin)
		}
	}
}
