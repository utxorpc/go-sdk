package sdk

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
)

type UtxorpcClient struct {
	httpClient     connect.HTTPClient
	baseUrl        string
	headers        map[string]string
	dialTimeout    time.Duration
	requestTimeout time.Duration
	Query          QueryServiceClient
	Submit         SubmitServiceClient
	Sync           SyncServiceClient
	Watch          WatchServiceClient
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

// WithDialTimeout sets the timeout duration for establishing a connection to the UTxO RPC server.
// The timeout applies to the initial dial operation. If the connection cannot be established
// within the specified duration, the dial operation will fail with a timeout error.
//
// This setting does not apply if a custom HTTP client is provided via WithHttpClient.
//
// See also: WithRequestTimeout for setting the timeout for individual requests.
func WithDialTimeout(timeout time.Duration) ClientOption {
	return func(u *UtxorpcClient) {
		u.dialTimeout = timeout
	}
}

// WithRequestTimeout sets the timeout duration for individual requests made by the UtxorpcClient.
// The timeout applies to each request operation and will cause the request to fail if it
// exceeds the specified duration. If not set, requests will use the default timeout behavior.
//
// This setting does not apply if a custom HTTP client is provided via WithHttpClient.
//
// See also: WithDialTimeout for setting the timeout for establishing connections.
func WithRequestTimeout(timeout time.Duration) ClientOption {
	return func(u *UtxorpcClient) {
		u.requestTimeout = timeout
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
			u.httpClient = createHttpClient(false, u.dialTimeout, u.requestTimeout)
		} else {
			u.httpClient = createHttpClient(true, u.dialTimeout, u.requestTimeout)
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

func createHttpClient(enableTls bool, dialTimeout, requestTimeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: requestTimeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, tlsConfig *tls.Config) (net.Conn, error) {
				if enableTls {
					// Establish a TLS connection using the custom TLS configuration
					conn, err := tls.DialWithDialer(&net.Dialer{Timeout: dialTimeout}, network, addr, tlsConfig)
					if err != nil {
						return nil, fmt.Errorf(
							"failed to establish TLS connection: %w",
							err,
						)
					}
					return conn, nil
				}
				return net.DialTimeout(network, addr, dialTimeout)
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
