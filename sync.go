package sdk

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

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
	// We assume blockIndex can be 0 or any positive number, and won't overflow
	// #nosec G115
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

func (u *UtxorpcClient) ReadTip() (*connect.Response[sync.ReadTipResponse], error) {
	return u.ReadTipWithContext(context.Background())
}

func (u *UtxorpcClient) ReadTipWithContext(
	ctx context.Context,
) (*connect.Response[sync.ReadTipResponse], error) {
	readTipReqProto := &sync.ReadTipRequest{}
	reqReadTip := connect.NewRequest(readTipReqProto)
	u.AddHeadersToRequest(reqReadTip)

	tipResp, err := u.Sync.ReadTip(ctx, reqReadTip)
	if err != nil {
		return nil, fmt.Errorf("failed to read tip: %w", err)
	}
	if tipResp.Msg == nil || tipResp.Msg.GetTip() == nil {
		return nil, errors.New("received nil tip from ReadTipResponse")
	}

	return tipResp, nil
}

func (u *UtxorpcClient) ReadBlock(
	blockRef *sync.BlockRef,
) (*connect.Response[sync.FetchBlockResponse], error) {
	return u.ReadBlockWithContext(context.Background(), blockRef)
}

func (u *UtxorpcClient) ReadBlockWithContext(
	ctx context.Context,
	blockRef *sync.BlockRef,
) (*connect.Response[sync.FetchBlockResponse], error) {
	fetchBlockReqProto := &sync.FetchBlockRequest{Ref: []*sync.BlockRef{blockRef}}
	reqFetchBlock := connect.NewRequest(fetchBlockReqProto)
	u.AddHeadersToRequest(reqFetchBlock)

	blockRespFull, err := u.Sync.FetchBlock(ctx, reqFetchBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block for tip: %w", err)
	}
	if blockRespFull.Msg == nil || len(blockRespFull.Msg.GetBlock()) == 0 || blockRespFull.Msg.GetBlock()[0] == nil {
		return nil, errors.New("received nil or empty block data from FetchBlockResponse for tip")
	}

	anyChainBlock := blockRespFull.Msg.GetBlock()[0]

	switch chain := anyChainBlock.GetChain().(type) {
	case *sync.AnyChainBlock_Cardano:
		if chain.Cardano != nil && chain.Cardano.GetHeader() != nil {
			return blockRespFull, nil
		} else {
			return nil, errors.New("cardano block or header is nil in FetchBlock response for tip")
		}
	default:
		return nil, fmt.Errorf("unknown or unsupported chain type in FetchBlock response: %T", chain)
	}
}
