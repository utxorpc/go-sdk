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
	var mode string = "submitTx"

	switch mode {
	case "submitTx":
		// Submit a transaction
		txCbor := "Replace this with signed CBOR transaction"
		txRefs, err := submitTx(ctx, client, txCbor)
		if err != nil {
			fmt.Printf("Error submitting transaction: %v\n", err)
			return
		}
		// Immediately wait for the transaction confirmation
		if err := waitForTx(ctx, client, txRefs); err != nil {
			fmt.Printf("Error waiting for transaction: %v\n", err)
		}
	case "readMempool":
		readMempool(ctx, client)
	case "waitForTx":
		if err := waitForTx(ctx, client, []string{"31bffedd962f4a6f5e85620985ccdf71f7b78988a6483e090f42d1e8badcebc8"}); err != nil {
			fmt.Printf("Error waiting for transaction: %v\n", err)
		}
	case "watchMempool":
		watchMempool(ctx, client)
	default:
		fmt.Println("Unknown mode:", mode)
	}
}

// Modified submitTx to return transaction references
func submitTx(ctx context.Context, client *utxorpc.UtxorpcClient, txCbor string) ([]string, error) {
	// Decode the transaction data from hex
	txRawBytes, err := hex.DecodeString(txCbor)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction hash: %w", err)
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
		if connectErr, ok := err.(*connect.Error); ok {
			// Extract error details
			errorCode := connectErr.Code()
			errorMessage := connectErr.Error()
			grpcMessage := connectErr.Meta().Get("Grpc-Message")
			return nil, fmt.Errorf(
				"gRPC error occurred:\n  Code: %v\n  Message: %s\n  Details: %s",
				errorCode,
				errorMessage,
				grpcMessage,
			)
		}
		return nil, fmt.Errorf("unexpected error occurred: %w", err)
	}

	// Extract and return transaction references
	if resp != nil && resp.Msg.Ref != nil {
		var refs []string
		fmt.Println("Response:")
		for i, ref := range resp.Msg.Ref {
			hexRef := hex.EncodeToString(ref)
			refs = append(refs, hexRef)
			fmt.Printf("  Ref[%d]: %s\n", i, hexRef)
		}
		return refs, nil
	}

	fmt.Println("No references found in the response.")
	return nil, nil
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

func waitForTx(ctx context.Context, client *utxorpc.UtxorpcClient, txRefs []string) error {
	fmt.Println("Waiting for the following transaction references:")
	for _, ref := range txRefs {
		fmt.Printf("  TxRef: %s\n", ref)
	}

	// Decode the transaction references from hex
	var decodedRefs [][]byte
	for _, ref := range txRefs {
		refBytes, err := hex.DecodeString(ref)
		if err != nil {
			return fmt.Errorf("failed to decode transaction reference %s: %w", ref, err)
		}
		decodedRefs = append(decodedRefs, refBytes)
	}

	// Create a WaitForTxRequest with the decoded transaction references
	req := connect.NewRequest(&submit.WaitForTxRequest{
		Ref: decodedRefs,
	})
	client.AddHeadersToRequest(req)

	fmt.Println("Connecting to utxorpc host:", client.URL())
	// Open a streaming connection to wait for transaction confirmation
	stream, err := client.Submit.WaitForTx(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to open waitForTx stream: %w", err)
	}
	defer stream.Close()

	// Process the stream of responses
	for stream.Receive() {
		resp := stream.Msg()

		// Decode and print the received stage and reference
		txRef := hex.EncodeToString(resp.Ref)
		txStage := resp.Stage
		fmt.Printf("Transaction %s is at stage: %v\n", txRef, txStage)

		// Break if the desired stage is reached (e.g., confirmed)
		if txStage == submit.Stage_STAGE_CONFIRMED {
			fmt.Printf("Transaction %s has been confirmed.\n", txRef)
			break
		}
	}

	// Check for stream errors
	if err := stream.Err(); err != nil {
		return fmt.Errorf("stream error: %w", err)
	}
	return nil
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
