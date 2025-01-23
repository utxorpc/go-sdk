package sdk

import (
	"context"
	"encoding/hex"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync/syncconnect"
)

type SyncServiceClient syncconnect.SyncServiceClient

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

func syncIntersect(blockHashStr string, blockIndex int64) []*sync.BlockRef {
	var intersect []*sync.BlockRef
	// Construct the BlockRef based on the provided parameters
	blockRef := &sync.BlockRef{}
	if blockHashStr != "" {
		hash, err := hex.DecodeString(blockHashStr)
		if err != nil {
			return nil
		}
		blockRef.Hash = hash
	}
	// We assume blockIndex can be 0 or any positive number
	if blockIndex > -1 {
		blockRef.Index = uint64(blockIndex)
	}
	// Only add blockRef to intersect if at least one of blockHashStr or blockIndex is provided
	if blockHashStr != "" || blockIndex > -1 {
		intersect = []*sync.BlockRef{blockRef}
	}
	return intersect
}

func (u *UtxorpcClient) FetchBlock(
	blockHashStr string,
	blockIndex int64,
) (*connect.Response[sync.FetchBlockResponse], error) {
	ctx := context.Background()
	req := &sync.FetchBlockRequest{Ref: syncIntersect(blockHashStr, blockIndex)}
	return u.FetchBlockWithContext(ctx, req)
}

func (u *UtxorpcClient) FetchBlockWithContext(
	ctx context.Context,
	blockReq *sync.FetchBlockRequest,
) (*connect.Response[sync.FetchBlockResponse], error) {
	req := connect.NewRequest(blockReq)
	u.AddHeadersToRequest(req)
	return u.Sync.FetchBlock(ctx, req)
}

func (u *UtxorpcClient) FollowTip(
	blockHashStr string,
	blockIndex int64,
) (*connect.ServerStreamForClient[sync.FollowTipResponse], error) {
	ctx := context.Background()
	req := &sync.FollowTipRequest{Intersect: syncIntersect(blockHashStr, blockIndex)}
	return u.FollowTipWithContext(ctx, req)
}

func (u *UtxorpcClient) FollowTipWithContext(
	ctx context.Context,
	blockReq *sync.FollowTipRequest,
) (*connect.ServerStreamForClient[sync.FollowTipResponse], error) {
	req := connect.NewRequest(blockReq)
	u.AddHeadersToRequest(req)
	return u.Sync.FollowTip(ctx, req)
}
