package sdk

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
)

type UtxorpcClient struct {
	httpClient connect.HTTPClient
	baseUrl    string
	headers    map[string]string
	Query      QueryServiceClient
	Submit     SubmitServiceClient
	Sync       SyncServiceClient
	Watch      WatchServiceClient
}

type ClientOption func(*UtxorpcClient)

func WithBaseUrl(baseUrl string) ClientOption {
	return func(u *UtxorpcClient) {
		u.baseUrl = baseUrl
	}
}

func WithHeaders(headers map[string]string) ClientOption {
	return func(u *UtxorpcClient) {
		u.headers = headers
	}
}

func WithHttpClient(httpClient connect.HTTPClient) ClientOption {
	return func(u *UtxorpcClient) {
		u.httpClient = httpClient
	}
}

func NewClient(options ...ClientOption) *UtxorpcClient {
	u := &UtxorpcClient{}

	for _, option := range options {
		option(u)
	}
	if u.httpClient == nil {
		if strings.HasPrefix(u.baseUrl, "http://") {
			u.httpClient = createHttpClient(false)
		} else {
			u.httpClient = createHttpClient(true)
		}
	}
	u.Query = u.NewQueryServiceClient()
	u.Submit = u.NewSubmitServiceClient()
	u.Sync = u.NewSyncServiceClient()
	u.Watch = u.NewWatchServiceClient()
	return u
}

func (u *UtxorpcClient) reset() {
	u.Query = u.NewQueryServiceClient()
	u.Submit = u.NewSubmitServiceClient()
	u.Sync = u.NewSyncServiceClient()
	u.Watch = u.NewWatchServiceClient()
}

func (u *UtxorpcClient) HTTPClient() connect.HTTPClient {
	return u.httpClient
}

func (u *UtxorpcClient) SetURL(baseUrl string) {
	u.baseUrl = baseUrl
	u.reset()
}

func (u *UtxorpcClient) URL() string {
	return u.baseUrl
}

func createHttpClient(enableTls bool) *http.Client {
	return &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, tlsConfig *tls.Config) (net.Conn, error) {
				if enableTls {
					// Establish a TLS connection using the custom TLS configuration
					conn, err := tls.Dial(network, addr, tlsConfig)
					if err != nil {
						return nil, fmt.Errorf(
							"failed to establish TLS connection: %w",
							err,
						)
					}
					return conn, nil
				}
				return net.Dial(network, addr)
			},
		},
	}
}

func (u *UtxorpcClient) Headers() map[string]string {
	headers := u.headers
	if headers == nil {
		headers = make(map[string]string)
	}
	return headers
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
