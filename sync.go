package sdk

import (
	"context"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1beta/sync"
	"github.com/utxorpc/go-codegen/utxorpc/v1beta/sync/syncconnect"
)

// SyncServiceClient is the generated Connect client for the UTxO RPC Sync
// service. Re-exported so callers can hold a typed reference without
// importing syncconnect directly.
type SyncServiceClient = syncconnect.SyncServiceClient

// NewSyncServiceClient is the free-function form of
// [(*UtxorpcClient).NewSyncServiceClient].
func NewSyncServiceClient(u *UtxorpcClient) SyncServiceClient {
	return u.NewSyncServiceClient()
}

// NewSyncServiceClient returns a fresh [SyncServiceClient] bound to this
// client's HTTP client and base URL.
func (u *UtxorpcClient) NewSyncServiceClient() SyncServiceClient {
	return syncconnect.NewSyncServiceClient(
		u.httpClient,
		u.baseUrl,
		connect.WithGRPC(),
	)
}

// FetchBlock calls [(*UtxorpcClient).FetchBlockWithContext] with a background context.
func (u *UtxorpcClient) FetchBlock(
	req *connect.Request[sync.FetchBlockRequest],
) (*connect.Response[sync.FetchBlockResponse], error) {
	ctx := context.Background()
	return u.FetchBlockWithContext(ctx, req)
}

// FetchBlockWithContext invokes Sync.FetchBlock after injecting stored
// headers into the request. Returns one or more blocks identified by their
// [sync.BlockRef] (hash + slot).
func (u *UtxorpcClient) FetchBlockWithContext(
	ctx context.Context,
	req *connect.Request[sync.FetchBlockRequest],
) (*connect.Response[sync.FetchBlockResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Sync.FetchBlock(ctx, req)
}

// FollowTip calls [(*UtxorpcClient).FollowTipWithContext] with a background context.
func (u *UtxorpcClient) FollowTip(
	req *connect.Request[sync.FollowTipRequest],
) (*connect.ServerStreamForClient[sync.FollowTipResponse], error) {
	ctx := context.Background()
	return u.FollowTipWithContext(ctx, req)
}

// FollowTipWithContext opens a server stream of chain-tip events: Apply for
// new blocks, Undo for rollbacks, and Reset when the server tells the client
// to restart from a given point. Callers MUST handle all three actions to
// maintain a consistent view of chain state. The stream must be closed.
func (u *UtxorpcClient) FollowTipWithContext(
	ctx context.Context,
	req *connect.Request[sync.FollowTipRequest],
) (*connect.ServerStreamForClient[sync.FollowTipResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Sync.FollowTip(ctx, req)
}

// ReadTip calls [(*UtxorpcClient).ReadTipWithContext] with a background context.
func (u *UtxorpcClient) ReadTip(
	req *connect.Request[sync.ReadTipRequest],
) (*connect.Response[sync.ReadTipResponse], error) {
	ctx := context.Background()
	return u.ReadTipWithContext(ctx, req)
}

// ReadTipWithContext invokes Sync.ReadTip after injecting stored headers
// into the request. Returns the current chain tip [sync.BlockRef].
func (u *UtxorpcClient) ReadTipWithContext(
	ctx context.Context,
	req *connect.Request[sync.ReadTipRequest],
) (*connect.Response[sync.ReadTipResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Sync.ReadTip(ctx, req)
}
