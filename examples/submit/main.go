package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/submit"
	utxorpc "github.com/utxorpc/go-sdk"
)

func main() {
	baseUrl := os.Getenv("UTXORPC_URL")
	if baseUrl == "" {
		baseUrl = "https://preview.utxorpc-v0.demeter.run"
	}
	client := utxorpc.NewClient(utxorpc.WithBaseUrl(baseUrl))
	dmtrApiKey := os.Getenv("DMTR_API_KEY")
	// set API key for demeter
	if dmtrApiKey != "" {
		client.SetHeader("dmtr-api-key", "dmtr_apikey...")
	}

	// Set mode to "submitTx", "readMempool", "waitForTx", or "watchMempool" to select the desired example.
	var mode string = "readMempool"

	switch mode {
	case "submitTx":
		// Submit a transaction
		txCbor := "Replace this with signed CBOR transaction"
		txRef, err := submitTx(client, txCbor)
		if err != nil {
			fmt.Printf("Error submitting transaction: %v\n", err)
			return
		}
		// Immediately wait for the transaction confirmation
		if err := waitForTx(client, txRef); err != nil {
			fmt.Printf("Error waiting for transaction: %v\n", err)
		}
	case "readMempool":
		readMempool(client)
	case "waitForTx":
		if err := waitForTx(client, "31bffedd962f4a6f5e85620985ccdf71f7b78988a6483e090f42d1e8badcebc8"); err != nil {
			fmt.Printf("Error waiting for transaction: %v\n", err)
		}
	case "watchMempool":
		watchMempool(client)
	default:
		fmt.Println("Unknown mode:", mode)
	}
}

// Modified submitTx to return transaction references
func submitTx(client *utxorpc.UtxorpcClient, txCbor string) (string, error) {
	fmt.Println("Connecting to utxorpc host:", client.URL())
	resp, err := client.SubmitTx(txCbor)
	if err != nil {
		if connectErr, ok := err.(*connect.Error); ok {
			// Extract error details
			errorCode := connectErr.Code()
			errorMessage := connectErr.Error()
			grpcMessage := connectErr.Meta().Get("Grpc-Message")
			return "", fmt.Errorf(
				"gRPC error occurred:\n  Code: %v\n  Message: %s\n  Details: %s",
				errorCode,
				errorMessage,
				grpcMessage,
			)
		}
		return "", fmt.Errorf("unexpected error occurred: %w", err)
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
		return refs[0], nil
	}

	return "", fmt.Errorf("No references found in the response.")
}

func readMempool(client *utxorpc.UtxorpcClient) {
	resp, err := client.ReadMempool()
	if err != nil {
		utxorpc.HandleError(err)
	}
	fmt.Printf("Response: %+v\n", resp)
}

func waitForTx(
	client *utxorpc.UtxorpcClient,
	txRef string,
) error {
	fmt.Println("Waiting for the following transaction reference:")
	fmt.Printf("  TxRef: %s\n", txRef)

	fmt.Println("Connecting to utxorpc host:", client.URL())
	// Open a streaming connection to wait for transaction confirmation
	stream, err := client.WaitForTx(txRef)
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

func watchMempool(client *utxorpc.UtxorpcClient) {
	fmt.Println("Connecting to utxorpc host:", client.URL())
	stream, err := client.WatchMempool()
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
