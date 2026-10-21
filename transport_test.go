package sdk_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	sdk "github.com/utxorpc/go-sdk"
	"github.com/utxorpc/go-sdk/v1alpha"
)

func TestDefaultHTTP2Transport(t *testing.T) {
	const requestTimeout = 5 * time.Second
	clients := []struct {
		name string
		new  func(string) connect.HTTPClient
	}{
		{
			name: "v1beta",
			new: func(url string) connect.HTTPClient {
				return sdk.NewClient(
					sdk.WithBaseUrl(url),
					sdk.WithDialTimeout(time.Second),
					sdk.WithRequestTimeout(requestTimeout),
				).HTTPClient()
			},
		},
		{
			name: "v1alpha",
			new: func(url string) connect.HTTPClient {
				return v1alpha.NewClient(
					v1alpha.WithBaseUrl(url),
					v1alpha.WithDialTimeout(time.Second),
					v1alpha.WithRequestTimeout(requestTimeout),
				).HTTPClient()
			},
		},
	}

	for _, client := range clients {
		t.Run(client.name, func(t *testing.T) {
			for _, useTLS := range []bool{false, true} {
				name := "h2c"
				if useTLS {
					name = "TLS"
				}
				t.Run(name, func(t *testing.T) {
					server := httptest.NewUnstartedServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if r.ProtoMajor != 2 {
								t.Errorf("request protocol = %s, want HTTP/2", r.Proto)
							}
							http.Redirect(w, r, "/redirected", http.StatusFound)
						}),
					)
					server.EnableHTTP2 = true
					if useTLS {
						server.StartTLS()
					} else {
						server.Config.Protocols = new(http.Protocols)
						server.Config.Protocols.SetUnencryptedHTTP2(true)
						server.Start()
					}
					defer server.Close()

					httpClient := client.new(server.URL).(*http.Client)
					defer httpClient.CloseIdleConnections()
					if httpClient.Timeout != requestTimeout {
						t.Fatalf("request timeout = %s, want %s",
							httpClient.Timeout, requestTimeout)
					}
					transport := httpClient.Transport.(*http.Transport)
					if useTLS {
						transport.TLSClientConfig.RootCAs = server.Client().
							Transport.(*http.Transport).TLSClientConfig.RootCAs
					}

					req, err := http.NewRequestWithContext(
						t.Context(), http.MethodGet, server.URL, nil,
					)
					if err != nil {
						t.Fatal(err)
					}
					resp, err := httpClient.Do(req)
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()
					if resp.ProtoMajor != 2 {
						t.Errorf("response protocol = %s, want HTTP/2", resp.Proto)
					}
					if resp.StatusCode != http.StatusFound {
						t.Errorf("status = %d, want unfollowed redirect", resp.StatusCode)
					}
				})
			}
		})
	}
}
