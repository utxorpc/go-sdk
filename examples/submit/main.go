package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/submit"
	utxorpc "github.com/utxorpc/go-sdk"
)

func main() {
	ctx := context.Background()
	baseUrl := "https://preview.utxorpc-v0.demeter.run"
	// set API key for demeter
	client := utxorpc.CreateUtxoRPCClient(baseUrl,
		// set API key for demeter
		utxorpc.WithHeaders(map[string]string{
			"dmtr-api-key": "dmtr_utxorpc1...",
		}),
	)

	// Set mode to "submitTx", "readMempool", "waitForTx", or "watchMempool" to select the desired example.
	var mode string = "waitForTx"

	switch mode {
	case "submitTx":
		submitTx(ctx, client, "Replace this with the signed transaction in CBOR format.")
	case "readMempool":
		readMempool(ctx, client)
	case "waitForTx":
		waitForTx(ctx, client)
	case "watchMempool":
		watchMempool(ctx, client)
	default:
		fmt.Println("Unknown mode:", mode)
	}
}

func submitTx(ctx context.Context, client *utxorpc.UtxorpcClient, txCbor string) {
	// Decode the transaction data from hex
	txRawBytes, err := hex.DecodeString(txCbor)
	if err != nil {
		panic(fmt.Errorf("failed to decode transaction hash: %v", err))
	}

	// Create a SubmitTxRequest with the transaction data
	tx := &submit.AnyChainTx{
		Type: &submit.AnyChainTx_Raw{
			Raw: txRawBytes,
		},
	}

	// Create a list with one transaction
	req := connect.NewRequest(&submit.SubmitTxRequest{
		Tx: []*submit.AnyChainTx{tx},
	})
	client.AddHeadersToRequest(req)

	fmt.Println("Connecting to utxorpc host:", client.URL())
	resp, err := client.Submit.SubmitTx(ctx, req)
	if err != nil {
		utxorpc.HandleError(err)
	}
	fmt.Printf("Response: %+v\n", resp)
}

func readMempool(ctx context.Context, client *utxorpc.UtxorpcClient) {
	req := connect.NewRequest(&submit.ReadMempoolRequest{})
	client.AddHeadersToRequest(req)
	fmt.Println("Connecting to utxorpc host:", client.URL())
	resp, err := client.Submit.ReadMempool(ctx, req)
	if err != nil {
		utxorpc.HandleError(err)
	}
	fmt.Printf("Response: %+v\n", resp)
}

func waitForTx(ctx context.Context, client *utxorpc.UtxorpcClient) {
	req := connect.NewRequest(&submit.WaitForTxRequest{})
	client.AddHeadersToRequest(req)
	fmt.Println("Connecting to utxorpc host:", client.URL())
	stream, err := client.Submit.WaitForTx(ctx, req)
	if err != nil {
		utxorpc.HandleError(err)
	}

	fmt.Println("Connected to utxorpc host, watching mempool...")
	for stream.Receive() {
		resp := stream.Msg()
		fmt.Printf("Stream response: %+v\n", resp)
	}

	if err := stream.Err(); err != nil {
		fmt.Println("Stream ended with error:", err)
	} else {
		fmt.Println("Stream ended normally.")
	}
}

func watchMempool(ctx context.Context, client *utxorpc.UtxorpcClient) {
	req := connect.NewRequest(&submit.WatchMempoolRequest{})
	client.AddHeadersToRequest(req)
	fmt.Println("Connecting to utxorpc host:", client.URL())
	stream, err := client.Submit.WatchMempool(ctx, req)
	if err != nil {
		utxorpc.HandleError(err)
	}

	fmt.Println("Connected to utxorpc host, watching mempool...")
	for stream.Receive() {
		resp := stream.Msg()
		fmt.Printf("Stream response: %+v\n", resp)
	}

	if err := stream.Err(); err != nil {
		fmt.Println("Stream ended with error:", err)
	} else {
		fmt.Println("Stream ended normally.")
	}
}
