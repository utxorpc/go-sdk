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
	followTip(ctx, client, "230eeba5de6b0198f64a3e801f92fa1ebf0f3a42a74dbd1922187249ad3038e7")
	followTip(ctx, client, "")
}

func fetchBlock(ctx context.Context, client *utxorpc.UtxorpcClient) {
	req := connect.NewRequest(&sync.FetchBlockRequest{})
	client.AddHeadersToRequest(req)

	fmt.Println("connecting to utxorpc host:", client.URL())
	chainSync, err := client.ChainSync.FetchBlock(ctx, req)
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
	resp, err := client.ChainSync.FollowTip(ctx, req)
	if err != nil {
		utxorpc.HandleError(err)
	}
	fmt.Println("connected to utxorpc...")
	fmt.Printf("Response: %+v\n", resp)
}
