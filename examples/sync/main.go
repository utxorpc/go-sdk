package main

import (
	"encoding/hex"
	"fmt"
	"os"

	sync "github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync"
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

	// Run them all
	fetchBlock(
		client,
		"235f9a217b826276d6cdfbb05c11572a06aef092535b6df8c682d501af59c230",
		65017558,
	)
	followTip(
		client,
		"235f9a217b826276d6cdfbb05c11572a06aef092535b6df8c682d501af59c230",
		65017558,
	)
}

func fetchBlock(
	client *utxorpc.UtxorpcClient,
	blockHash string,
	blockIndex int64,
) {
	fmt.Println("connecting to utxorpc host:", client.URL())
	resp, err := client.FetchBlock(blockHash, blockIndex)
	if err != nil {
		utxorpc.HandleError(err)
	}
	fmt.Printf("Response: %+v\n", resp)
	for i, blockRef := range resp.Msg.Block {
		fmt.Printf("Block[%d]:\n", i)
		fmt.Printf("Index: %d\n", blockRef.GetCardano().GetHeader().GetSlot())
		fmt.Printf("Hash: %x\n", blockRef.GetCardano().GetHeader().GetHash())
	}
}

func followTip(
	client *utxorpc.UtxorpcClient,
	blockHash string,
	blockIndex int64,
) {
	fmt.Println("connecting to utxorpc host:", client.URL())
	stream, err := client.FollowTip(blockHash, blockIndex)
	if err != nil {
		utxorpc.HandleError(err)
		return
	}
	fmt.Println("Connected to utxorpc host, following tip...")

	for stream.Receive() {
		resp := stream.Msg()
		action := resp.GetAction()
		switch a := action.(type) {
		case *sync.FollowTipResponse_Apply:
			fmt.Println("Action: Apply")
			printAnyChainBlock(a.Apply)
		case *sync.FollowTipResponse_Undo:
			fmt.Println("Action: Undo")
			printAnyChainBlock(a.Undo)
		case *sync.FollowTipResponse_Reset_:
			fmt.Println("Action: Reset")
			printBlockRef(a.Reset_)
		default:
			fmt.Println("Unknown action type")
		}
	}

	if err := stream.Err(); err != nil {
		fmt.Println("Stream ended with error:", err)
	} else {
		fmt.Println("Stream ended normally.")
	}
}

func printAnyChainBlock(block *sync.AnyChainBlock) {
	if block == nil {
		return
	}
	if cardanoBlock := block.GetCardano(); cardanoBlock != nil {
		hash := hex.EncodeToString(cardanoBlock.Header.Hash)
		slot := cardanoBlock.Header.Slot
		fmt.Printf("Block Slot: %d, Block Hash: %s\n", slot, hash)
	}
}

func printBlockRef(blockRef *sync.BlockRef) {
	if blockRef == nil {
		return
	}
	hash := hex.EncodeToString(blockRef.Hash)
	slot := blockRef.Index
	fmt.Printf("Block Slot: %d, Block Hash: %s\n", slot, hash)
}
