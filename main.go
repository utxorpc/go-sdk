package sdk

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/query/queryconnect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/submit/submitconnect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync/syncconnect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/watch/watchconnect"
	"golang.org/x/net/http2"
)

type UtxorpcClient struct {
	httpClient connect.HTTPClient
	baseUrl    string
	apiKey     string
	ChainSync  syncconnect.ChainSyncServiceClient
	Query      queryconnect.QueryServiceClient
	Submit     submitconnect.SubmitServiceClient
	Watch      watchconnect.WatchServiceClient
}

func NewClient(httpClient *http.Client, baseUrl string, apiKey string) UtxorpcClient {
	var client UtxorpcClient
	chainSyncClient := syncconnect.NewChainSyncServiceClient(httpClient, baseUrl, connect.WithGRPC())
	queryClient := queryconnect.NewQueryServiceClient(httpClient, baseUrl, connect.WithGRPC())
	submitClient := submitconnect.NewSubmitServiceClient(httpClient, baseUrl, connect.WithGRPC())
	watchClient := watchconnect.NewWatchServiceClient(httpClient, baseUrl, connect.WithGRPC())
	client.httpClient = httpClient
	client.baseUrl = baseUrl
	client.apiKey = apiKey
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

func CreateUtxoRPCClient(baseUrl, apiKey string) *UtxorpcClient {
	httpClient := createHttpClient()
	client := NewClient(httpClient, baseUrl, apiKey)
	return &client
}

func (u *UtxorpcClient) SetAPIKeyHeader(req connect.AnyRequest) {
	req.Header().Set("dmtr-api-key", u.apiKey)
}

func HandleError(err error) {
	fmt.Println(connect.CodeOf(err))
	if connectErr := new(connect.Error); errors.As(err, &connectErr) {
		fmt.Println(connectErr.Message())
		fmt.Println(connectErr.Details())
	}
	panic(err)
}
