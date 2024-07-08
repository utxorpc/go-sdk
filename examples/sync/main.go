package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"connectrpc.com/connect"
	sync "github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync"
	utxorpc "github.com/utxorpc/go-sdk"
)

func main() {
	ctx := context.Background()
	baseUrl := "https://preview.utxorpc-v0.demeter.run"
	client := utxorpc.CreateUtxoRPCClient(baseUrl,
		// set API key for demeter
		utxorpc.WithHeaders(map[string]string{
			"dmtr-api-key": "dmtr_utxorpc1...",
		}),
	)

	fetchBlock(ctx, client)
	followTip(ctx, client, "5d316b4722f832bb7cee7914cbcc7b977a03fea6884da5f2157a6ae420b91280")
	followTip(ctx, client, "")
}

func fetchBlock(ctx context.Context, client *utxorpc.UtxorpcClient) {
	req := connect.NewRequest(&sync.FetchBlockRequest{})
	client.AddHeadersToRequest(req)

	fmt.Println("connecting to utxorpc host:", client.URL())
	chainSync, err := client.Sync.FetchBlock(ctx, req)
	if err != nil {
		utxorpc.HandleError(err)
	}
	fmt.Println("connected to utxorpc...")
	for i, blockRef := range chainSync.Msg.Block {
		fmt.Printf("Block[%d]:\n", i)
		fmt.Printf("Index: %d\n", blockRef.GetCardano().GetHeader().GetSlot())
		fmt.Printf("Hash: %x\n", blockRef.GetCardano().GetHeader().GetHash())
	}
}

// FollowTipRequest with Intersect
func followTip(ctx context.Context, client *utxorpc.UtxorpcClient, blockHash string) {
	var req *connect.Request[sync.FollowTipRequest]

	if blockHash == "" {
		req = connect.NewRequest(&sync.FollowTipRequest{})
	} else {
		hash, err := hex.DecodeString(blockHash)
		if err != nil {
			log.Fatalf("failed to decode hex string: %v", err)
		}

		blockRef := &sync.BlockRef{
			Hash: hash,
		}
		req = connect.NewRequest(&sync.FollowTipRequest{
			Intersect: []*sync.BlockRef{blockRef},
		})
	}
	client.AddHeadersToRequest(req)
	fmt.Println("connecting to utxorpc host:", client.URL())
	stream, err := client.Sync.FollowTip(ctx, req)
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
