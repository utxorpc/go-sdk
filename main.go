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
	headers    map[string]string
	Query      queryconnect.QueryServiceClient
	Submit     submitconnect.SubmitServiceClient
	Sync       syncconnect.SyncServiceClient
	Watch      watchconnect.WatchServiceClient
}

type ClientOption func(*UtxorpcClient)

func WithHeaders(headers map[string]string) ClientOption {
	return func(client *UtxorpcClient) {
		client.headers = headers
	}
}

func NewClient(httpClient *http.Client, baseUrl string, options ...ClientOption) *UtxorpcClient {
	client := &UtxorpcClient{
		httpClient: httpClient,
		baseUrl:    baseUrl,
		Query:      queryconnect.NewQueryServiceClient(httpClient, baseUrl, connect.WithGRPC()),
		Submit:     submitconnect.NewSubmitServiceClient(httpClient, baseUrl, connect.WithGRPC()),
		Sync:       syncconnect.NewSyncServiceClient(httpClient, baseUrl, connect.WithGRPC()),
		Watch:      watchconnect.NewWatchServiceClient(httpClient, baseUrl, connect.WithGRPC()),
	}

	for _, option := range options {
		option(client)
	}
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

func CreateUtxoRPCClient(baseUrl string, options ...ClientOption) *UtxorpcClient {
	httpClient := createHttpClient()
	return NewClient(httpClient, baseUrl, options...)
}

func (u *UtxorpcClient) SetHeader(key, value string) {
	if u.headers == nil {
		u.headers = make(map[string]string)
	}
	u.headers[key] = value
}

func (u *UtxorpcClient) SetHeaders(headers map[string]string) {
	u.headers = headers
}

func (u *UtxorpcClient) RemoveHeader(key string) {
	delete(u.headers, key)
}

func (u *UtxorpcClient) AddHeadersToRequest(req connect.AnyRequest) {
	for key, value := range u.headers {
		req.Header().Set(key, value)
	}
}

func HandleError(err error) {
	fmt.Println(connect.CodeOf(err))
	if connectErr := new(connect.Error); errors.As(err, &connectErr) {
		fmt.Println(connectErr.Message())
		fmt.Println(connectErr.Details())
	}
	panic(err)
}
