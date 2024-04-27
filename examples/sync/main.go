package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"connectrpc.com/connect"
	sync "github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync"
	utxorpc "github.com/utxorpc/go-sdk"
	"golang.org/x/net/http2"
)

func main() {
	ctx := context.Background()
	httpClient := &http.Client{
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
	baseUrl := "https://preview.utxorpc-v0.demeter.run"
	client := utxorpc.NewClient(httpClient, baseUrl)
	req := connect.NewRequest(&sync.FetchBlockRequest{})
	// set API key for demeter
	req.Header().Set("dmtr-api-key", "dmtr_utxorpc1...")
	fmt.Println("connecting to utxorpc host:", baseUrl)
	chainSync, err := client.ChainSync.FetchBlock(ctx, req)
	if err != nil {
		fmt.Println(connect.CodeOf(err))
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			fmt.Println(connectErr.Message())
			fmt.Println(connectErr.Details())
		}
		panic(err)
		os.Exit(1)
	}
	fmt.Println("connected to utxorpc...")
	for i, blockRef := range chainSync.Msg.Block {
		fmt.Printf("Block[%d]:\n", i)
		fmt.Printf("Index: %d\n", blockRef.GetCardano().GetHeader().GetSlot())
		fmt.Printf("Hash: %x\n", blockRef.GetCardano().GetHeader().GetHash())
	}
}
