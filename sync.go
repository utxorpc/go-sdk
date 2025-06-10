package sdk

import (
	"context"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync/syncconnect"
)

type SyncServiceClient = syncconnect.SyncServiceClient

func NewSyncServiceClient(u *UtxorpcClient) SyncServiceClient {
	return u.NewSyncServiceClient()
}

func (u *UtxorpcClient) NewSyncServiceClient() SyncServiceClient {
	return syncconnect.NewSyncServiceClient(
		u.httpClient,
		u.baseUrl,
		connect.WithGRPC(),
	)
}

func (u *UtxorpcClient) FetchBlock(
	req *connect.Request[sync.FetchBlockRequest],
) (*connect.Response[sync.FetchBlockResponse], error) {
	ctx := context.Background()
	return u.FetchBlockWithContext(ctx, req)
}

func (u *UtxorpcClient) FetchBlockWithContext(
	ctx context.Context,
	req *connect.Request[sync.FetchBlockRequest],
) (*connect.Response[sync.FetchBlockResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Sync.FetchBlock(ctx, req)
}

func (u *UtxorpcClient) FollowTip(
	req *connect.Request[sync.FollowTipRequest],
) (*connect.ServerStreamForClient[sync.FollowTipResponse], error) {
	ctx := context.Background()
	return u.FollowTipWithContext(ctx, req)
}

func (u *UtxorpcClient) FollowTipWithContext(
	ctx context.Context,
	req *connect.Request[sync.FollowTipRequest],
) (*connect.ServerStreamForClient[sync.FollowTipResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Sync.FollowTip(ctx, req)
}

func (u *UtxorpcClient) ReadTip(
	req *connect.Request[sync.ReadTipRequest],
) (*connect.Response[sync.ReadTipResponse], error) {
	ctx := context.Background()
	return u.ReadTipWithContext(ctx, req)
}

func (u *UtxorpcClient) ReadTipWithContext(
	ctx context.Context,
	req *connect.Request[sync.ReadTipRequest],
) (*connect.Response[sync.ReadTipResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Sync.ReadTip(ctx, req)
}
