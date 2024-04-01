package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/build"
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
	baseUrl := "https://localhost:51000"
	client := utxorpc.NewClient(httpClient, baseUrl)
	req := connect.NewRequest(&build.GetChainTipRequest{})
	fmt.Println("connecting to utxorpc host:", baseUrl)
	chainTip, err := client.LedgerState.GetChainTip(ctx, req)
	if err != nil {
		fmt.Println(connect.CodeOf(err))
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			fmt.Println(connectErr.Message())
			fmt.Println(connectErr.Details())
		}
		panic(err)
	}
	fmt.Println("connected to utxorpc...")
	fmt.Printf("Chain Tip:\n")
	fmt.Printf("Slot: %d\n", chainTip.Msg.Tip.Slot)
	fmt.Printf("Height: %d\n", chainTip.Msg.Tip.Height)
	fmt.Printf("Hash: %x\n", chainTip.Msg.Tip.Hash)
}
