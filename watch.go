package sdk

import (
	"context"
	"encoding/hex"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/watch"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/watch/watchconnect"
)

type WatchServiceClient watchconnect.WatchServiceClient

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

func watchIntersect(blockHashStr string, blockIndex int64) []*watch.BlockRef {
	var intersect []*watch.BlockRef
	// Construct the BlockRef based on the provided parameters
	blockRef := &watch.BlockRef{}
	if blockHashStr != "" {
		hash, err := hex.DecodeString(blockHashStr)
		if err != nil {
			return nil
		}
		blockRef.Hash = hash
	}
	// We assume blockIndex can be 0 or any positive number, and won't overflow
	// #nosec G115
	if blockIndex > -1 {
		blockRef.Index = uint64(blockIndex)
	}
	// Only add blockRef to intersect if at least one of blockHashStr or blockIndex is provided
	if blockHashStr != "" || blockIndex > -1 {
		intersect = []*watch.BlockRef{blockRef}
	}
	return intersect
}

func (u *UtxorpcClient) WatchTx(
	blockHashStr string,
	blockIndex int64,
) (*connect.ServerStreamForClient[watch.WatchTxResponse], error) {
	ctx := context.Background()
	req := &watch.WatchTxRequest{Intersect: watchIntersect(blockHashStr, blockIndex)}
	return u.WatchTxWithContext(ctx, req)
}

func (u *UtxorpcClient) WatchTxWithContext(
	ctx context.Context,
	watchReq *watch.WatchTxRequest,
) (*connect.ServerStreamForClient[watch.WatchTxResponse], error) {
	req := connect.NewRequest(watchReq)
	u.AddHeadersToRequest(req)
	return u.Watch.WatchTx(ctx, req)
}
