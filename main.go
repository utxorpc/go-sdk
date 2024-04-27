package sdk

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/query/queryconnect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/submit/submitconnect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync/syncconnect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/watch/watchconnect"
)

type UtxorpcClient struct {
	httpClient connect.HTTPClient
	baseUrl    string
	ChainSync  syncconnect.ChainSyncServiceClient
	Query      queryconnect.QueryServiceClient
	Submit     submitconnect.SubmitServiceClient
	Watch      watchconnect.WatchServiceClient
}

func NewClient(httpClient *http.Client, baseUrl string) UtxorpcClient {
	var client UtxorpcClient
	chainSyncClient := syncconnect.NewChainSyncServiceClient(httpClient, baseUrl, connect.WithGRPC())
	queryClient := queryconnect.NewQueryServiceClient(httpClient, baseUrl, connect.WithGRPC())
	submitClient := submitconnect.NewSubmitServiceClient(httpClient, baseUrl, connect.WithGRPC())
	watchClient := watchconnect.NewWatchServiceClient(httpClient, baseUrl, connect.WithGRPC())
	client.httpClient = httpClient
	client.baseUrl = baseUrl
	client.ChainSync = chainSyncClient
	client.Query = queryClient
	client.Submit = submitClient
	client.Watch = watchClient
	return client
}

func (u *UtxorpcClient) HTTPClient() connect.HTTPClient {
	return u.httpClient
}

func (u *UtxorpcClient) URL() string {
	return u.baseUrl
}
