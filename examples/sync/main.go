package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"connectrpc.com/connect"
	sync "github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync"
	utxorpc "github.com/utxorpc/go-sdk"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
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

	// Set mode to "fetchBlock" or "followTip" to select the desired example.
	var mode string = "followTip"

	switch mode {
	case "fetchBlock":
		fetchBlock(ctx, client, "235f9a217b826276d6cdfbb05c11572a06aef092535b6df8c682d501af59c230", 65017558, nil)
	case "followTip":
		followTip(ctx, client, "235f9a217b826276d6cdfbb05c11572a06aef092535b6df8c682d501af59c230", 65017558, nil)
	default:
		fmt.Println("Unknown mode:", mode)
	}
}

func fetchBlock(ctx context.Context, client *utxorpc.UtxorpcClient, blockHash string, blockIndex int64, fieldMaskPaths []string) {
	var req *connect.Request[sync.FetchBlockRequest]
	var intersect []*sync.BlockRef
	var fieldMask *fieldmaskpb.FieldMask

	// Construct the BlockRef based on the provided parameters
	blockRef := &sync.BlockRef{}
	if blockHash != "" {
		hash, err := hex.DecodeString(blockHash)
		if err != nil {
			log.Fatalf("failed to decode hex string: %v", err)
		}
		blockRef.Hash = hash
	}
	// We assume blockIndex can be 0 or any positive number
	if blockIndex > -1 {
		blockRef.Index = uint64(blockIndex)
	}

	// Only add blockRef to intersect if at least one of blockHash or blockIndex is provided
	if blockHash != "" || blockIndex > -1 {
		intersect = []*sync.BlockRef{blockRef}
	}

	// Construct the FieldMask if paths are provided
	if len(fieldMaskPaths) > 0 {
		fieldMask = &fieldmaskpb.FieldMask{
			Paths: fieldMaskPaths,
		}
	}

	// Create the FetchBlockRequest
	req = connect.NewRequest(&sync.FetchBlockRequest{
		Ref: intersect,
		FieldMask: fieldMask,
	})

	// Print BlockRef details if intersect is provided
	if len(intersect) > 0 {
		fmt.Printf("Blockref: %d, %x\n", req.Msg.Ref[0].Index, req.Msg.Ref[0].Hash)
	}

	client.AddHeadersToRequest(req)
	fmt.Println("connecting to utxorpc host:", client.URL())
	resp, err := client.Sync.FetchBlock(ctx, req)
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

func followTip(ctx context.Context, client *utxorpc.UtxorpcClient, blockHash string, blockIndex int64, fieldMaskPaths []string) {
	var req *connect.Request[sync.FollowTipRequest]
	var intersect []*sync.BlockRef
	var fieldMask *fieldmaskpb.FieldMask

	// Construct the BlockRef based on the provided parameters
	blockRef := &sync.BlockRef{}
	if blockHash != "" {
		hash, err := hex.DecodeString(blockHash)
		if err != nil {
			log.Fatalf("failed to decode hex string: %v", err)
		}
		blockRef.Hash = hash
	}
	// We assume blockIndex can be 0 or any positive number
	if blockIndex > -1 {
		blockRef.Index = uint64(blockIndex)
	}

	// Only add blockRef to intersect if at least one of blockHash or blockIndex is provided
	if blockHash != "" || blockIndex > -1 {
		intersect = []*sync.BlockRef{blockRef}
	}

	// Construct the FieldMask if paths are provided
	if len(fieldMaskPaths) > 0 {
		fieldMask = &fieldmaskpb.FieldMask{
			Paths: fieldMaskPaths,
		}
	}

	// Create the FollowTipRequest
	req = connect.NewRequest(&sync.FollowTipRequest{
		Intersect: intersect,
		FieldMask: fieldMask,
	})

	// Print BlockRef details if intersect is provided
	if len(intersect) > 0 {
		fmt.Printf("Blockref: %d, %x\n", req.Msg.Intersect[0].Index, req.Msg.Intersect[0].Hash)
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
