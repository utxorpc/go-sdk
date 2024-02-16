package sdk

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync/syncconnect"
)

type UtxorpcClient struct {
	httpClient connect.HTTPClient
	baseUrl    string
	ChainSync  syncconnect.ChainSyncServiceClient
}

func NewClient(httpClient *http.Client, baseUrl string) UtxorpcClient {
	var client UtxorpcClient
	chainSyncClient := syncconnect.NewChainSyncServiceClient(httpClient, baseUrl, connect.WithGRPC())
	client.httpClient = httpClient
	client.baseUrl = baseUrl
	client.ChainSync = chainSyncClient
	return client
}

func (u *UtxorpcClient) URL() string {
	return u.baseUrl
}
