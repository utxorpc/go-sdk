package main

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"connectrpc.com/connect"
	sync "github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync"
	utxorpc "github.com/utxorpc/go-sdk"
	"golang.org/x/net/http2"
)

func main() {
	ctx := context.Background()
	baseUrl := "https://preview.utxorpc-v0.demeter.run"
	// set API key for demeter
	apiKey := "dmtr_utxorpc1..."
	client := createUtxoRPCClient(baseUrl)

	fetchBlock(ctx, client, apiKey)
	followTip(ctx, client, apiKey, "230eeba5de6b0198f64a3e801f92fa1ebf0f3a42a74dbd1922187249ad3038e7")
	followTip(ctx, client, apiKey, "")
}

func createHttpClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
				// If you're also using this client for non-h2c traffic, you may want
				// to delegate to tls.Dial if the network isn't TCP or the addr isn't
				// in an allowlist.
				return net.Dial(network, addr)
			},
		},
	}
}

func createUtxoRPCClient(baseUrl string) *utxorpc.UtxorpcClient {
	httpClient := createHttpClient()
	client := utxorpc.NewClient(httpClient, baseUrl)
	return &client
}

func setAPIKeyHeader(req connect.AnyRequest, apiKey string) {
	req.Header().Set("dmtr-api-key", apiKey)
}

func fetchBlock(ctx context.Context, client *utxorpc.UtxorpcClient, apiKey string) {
	req := connect.NewRequest(&sync.FetchBlockRequest{})
	setAPIKeyHeader(req, apiKey)

	fmt.Println("connecting to utxorpc host:", client.URL())
	chainSync, err := client.ChainSync.FetchBlock(ctx, req)
	if err != nil {
		handleError(err)
	}
	fmt.Println("connected to utxorpc...")
	for i, blockRef := range chainSync.Msg.Block {
		fmt.Printf("Block[%d]:\n", i)
		fmt.Printf("Index: %d\n", blockRef.GetCardano().GetHeader().GetSlot())
		fmt.Printf("Hash: %x\n", blockRef.GetCardano().GetHeader().GetHash())
	}
}

// FollowTipRequest with Intersect
func followTip(ctx context.Context, client *utxorpc.UtxorpcClient, apiKey string, blockHash string) {
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

	setAPIKeyHeader(req, apiKey)

	fmt.Println("connecting to utxorpc host:", client.URL())
	resp, err := client.ChainSync.FollowTip(ctx, req)
	if err != nil {
		handleError(err)
	}
	fmt.Println("connected to utxorpc...")
	fmt.Printf("Response: %+v\n", resp)
}

func handleError(err error) {
	fmt.Println(connect.CodeOf(err))
	if connectErr := new(connect.Error); errors.As(err, &connectErr) {
		fmt.Println(connectErr.Message())
		fmt.Println(connectErr.Details())
	}
	panic(err)
}
