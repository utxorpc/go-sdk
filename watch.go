package sdk

import (
	"context"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/watch"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/watch/watchconnect"
)

type WatchServiceClient = watchconnect.WatchServiceClient

func NewWatchServiceClient(u *UtxorpcClient) WatchServiceClient {
	return u.NewWatchServiceClient()
}

func (u *UtxorpcClient) NewWatchServiceClient() WatchServiceClient {
	return watchconnect.NewWatchServiceClient(
		u.httpClient,
		u.baseUrl,
		connect.WithGRPC(),
	)
}

func (u *UtxorpcClient) WatchTx(
	req *connect.Request[watch.WatchTxRequest],
) (*connect.ServerStreamForClient[watch.WatchTxResponse], error) {
	ctx := context.Background()
	return u.WatchTxWithContext(ctx, req)
}

func (u *UtxorpcClient) WatchTxWithContext(
	ctx context.Context,
	req *connect.Request[watch.WatchTxRequest],
) (*connect.ServerStreamForClient[watch.WatchTxResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Watch.WatchTx(ctx, req)
}
