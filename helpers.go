package sdk

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/cardano"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/query"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/submit"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/watch"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func (u *UtxorpcClient) EvaluateTransaction(
	txCbor string,
) (*connect.Response[submit.EvalTxResponse], error) {
	ctx := context.Background()
	// Decode the transaction data from hex
	txRawBytes, err := hex.DecodeString(txCbor)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction hash: %w", err)
	}

	// Create a EvalTxRequest with the transaction data
	tx := &submit.AnyChainTx{
		Type: &submit.AnyChainTx_Raw{
			Raw: txRawBytes,
		},
	}

	// Create a list with one transaction
	req := &submit.EvalTxRequest{
		Tx: []*submit.AnyChainTx{tx},
	}
	return u.EvaluateTransactionWithContext(ctx, req)
}

func (u *UtxorpcClient) EvaluateTransactionWithContext(
	ctx context.Context,
	txReq *submit.EvalTxRequest,
) (*connect.Response[submit.EvalTxResponse], error) {
	req := connect.NewRequest(txReq)
	u.AddHeadersToRequest(req)
	return u.Submit.EvalTx(ctx, req)
}

func (u *UtxorpcClient) GetMempoolTransactions() (*connect.Response[submit.ReadMempoolResponse], error) {
	ctx := context.Background()
	return u.GetMempoolTransactionsWithContext(ctx)
}

func (u *UtxorpcClient) GetMempoolTransactionsWithContext(
	ctx context.Context,
) (*connect.Response[submit.ReadMempoolResponse], error) {
	req := connect.NewRequest(&submit.ReadMempoolRequest{})
	u.AddHeadersToRequest(req)
	return u.Submit.ReadMempool(ctx, req)
}

func (u *UtxorpcClient) GetProtocolParameters() (*connect.Response[query.ReadParamsResponse], error) {
	ctx := context.Background()
	return u.GetProtocolParametersWithContext(ctx)
}

func (u *UtxorpcClient) GetProtocolParametersWithContext(
	ctx context.Context,
) (*connect.Response[query.ReadParamsResponse], error) {
	req := connect.NewRequest(&query.ReadParamsRequest{})
	u.AddHeadersToRequest(req)
	return u.Query.ReadParams(ctx, req)
}

func (u *UtxorpcClient) GetUtxoByRef(
	txHashStr string,
	txIndex uint32,
) (*connect.Response[query.ReadUtxosResponse], error) {
	return u.GetUtxoByRefWithContext(context.Background(), txHashStr, txIndex)
}

func (u *UtxorpcClient) GetUtxoByRefWithContext(
	ctx context.Context,
	txHashStr string,
	txIndex uint32,
) (*connect.Response[query.ReadUtxosResponse], error) {
	var txHashBytes []byte
	var err error
	// Attempt to decode the input as hex
	txHashBytes, hexErr := hex.DecodeString(txHashStr)
	if hexErr != nil {
		// If not hex, attempt to decode as Base64
		txHashBytes, err = base64.StdEncoding.DecodeString(txHashStr)
		if err != nil {
			return nil, err
		}
	}
	// Create TxoRef with the decoded hash bytes
	txoRef := &query.TxoRef{
		Hash:  txHashBytes, // Use the decoded []byte
		Index: txIndex,
	}
	txReq := &query.ReadUtxosRequest{Keys: []*query.TxoRef{txoRef}}
	req := connect.NewRequest(txReq)
	return u.ReadUtxosWithContext(ctx, req)
}

func (u *UtxorpcClient) GetUtxosByRefs(
	refs []*query.TxoRef,
) (*connect.Response[query.ReadUtxosResponse], error) {
	return u.GetUtxosByRefsWithContext(context.Background(), refs)
}

func (u *UtxorpcClient) GetUtxosByRefsWithContext(
	ctx context.Context,
	refs []*query.TxoRef,
) (*connect.Response[query.ReadUtxosResponse], error) {
	if len(refs) == 0 {
		return nil, errors.New("no transaction references provided")
	}

	txReq := &query.ReadUtxosRequest{Keys: refs}
	req := connect.NewRequest(txReq)
	return u.ReadUtxosWithContext(ctx, req)
}

func (u *UtxorpcClient) GetUtxosByAddress(
	address []byte,
) (*connect.Response[query.SearchUtxosResponse], error) {
	return u.GetUtxosByAddressWithContext(context.Background(), address)
}

func (u *UtxorpcClient) GetUtxosByAddressWithContext(
	ctx context.Context,
	address []byte,
) (*connect.Response[query.SearchUtxosResponse], error) {
	queryReq := &query.SearchUtxosRequest{
		FieldMask: &fieldmaskpb.FieldMask{Paths: []string{}},
		Predicate: &query.UtxoPredicate{
			Match: &query.AnyUtxoPattern{
				UtxoPattern: &query.AnyUtxoPattern_Cardano{
					Cardano: &cardano.TxOutputPattern{
						Address: &cardano.AddressPattern{
							ExactAddress: address,
						},
					},
				},
			},
		},
		MaxItems:   100, // May need adjustment
		StartToken: "",  // For pagination, start at first page
	}
	req := connect.NewRequest(queryReq)
	return u.SearchUtxosWithContext(ctx, req)
}

func (u *UtxorpcClient) GetUtxosByAddressWithAsset(
	addressBytes []byte,
	policyIdBytes []byte,
	assetNameBytes []byte,
) (*connect.Response[query.SearchUtxosResponse], error) {
	return u.GetUtxosByAddressWithAssetWithContext(
		context.Background(),
		addressBytes,
		policyIdBytes,
		assetNameBytes,
	)
}

func (u *UtxorpcClient) GetUtxosByAddressWithAssetWithContext(
	ctx context.Context,
	addressBytes []byte,
	policyIdBytes []byte,
	assetNameBytes []byte,
) (*connect.Response[query.SearchUtxosResponse], error) {
	tpl := &cardano.TxOutputPattern{
		Address: &cardano.AddressPattern{
			ExactAddress: addressBytes,
		},
	}

	var assetFilter *cardano.AssetPattern

	if len(policyIdBytes) > 0 && len(assetNameBytes) > 0 {
		assetFilter = &cardano.AssetPattern{
			PolicyId:  policyIdBytes,
			AssetName: assetNameBytes,
		}
	} else if len(policyIdBytes) > 0 {
		assetFilter = &cardano.AssetPattern{
			PolicyId: policyIdBytes,
		}
	} else if len(assetNameBytes) > 0 {
		assetFilter = &cardano.AssetPattern{
			AssetName: assetNameBytes,
		}
	}

	if assetFilter != nil {
		tpl.Asset = assetFilter
	}

	queryReq := &query.SearchUtxosRequest{
		FieldMask: &fieldmaskpb.FieldMask{Paths: []string{}},
		Predicate: &query.UtxoPredicate{
			Match: &query.AnyUtxoPattern{
				UtxoPattern: &query.AnyUtxoPattern_Cardano{
					Cardano: tpl,
				},
			},
		},
		MaxItems:   100, // May need adjustment
		StartToken: "",  // For pagination, start at first page
	}
	req := connect.NewRequest(queryReq)
	return u.SearchUtxosWithContext(ctx, req)
}

func (u *UtxorpcClient) GetUtxosByAsset(
	policyIdBytes []byte,
	assetNameBytes []byte,
) (*connect.Response[query.SearchUtxosResponse], error) {
	return u.GetUtxosByAssetWithContext(
		context.Background(),
		policyIdBytes,
		assetNameBytes,
	)
}

func (u *UtxorpcClient) GetUtxosByAssetWithContext(
	ctx context.Context,
	policyIdBytes []byte,
	assetNameBytes []byte,
) (*connect.Response[query.SearchUtxosResponse], error) {
	if policyIdBytes == nil && assetNameBytes == nil {
		return nil, errors.New(
			"at least one of policyId or assetName must be provided",
		)
	}

	assetPattern := &cardano.AssetPattern{}
	hasAssetFilter := false
	if policyIdBytes != nil {
		assetPattern.PolicyId = policyIdBytes
		hasAssetFilter = true
	}
	if assetNameBytes != nil {
		assetPattern.AssetName = assetNameBytes
		hasAssetFilter = true
	}

	cardanoOutputPattern := &cardano.TxOutputPattern{}

	if hasAssetFilter {
		cardanoOutputPattern.Asset = assetPattern
	}

	queryReq := &query.SearchUtxosRequest{
		FieldMask: &fieldmaskpb.FieldMask{Paths: []string{}},
		Predicate: &query.UtxoPredicate{
			Match: &query.AnyUtxoPattern{
				UtxoPattern: &query.AnyUtxoPattern_Cardano{
					Cardano: cardanoOutputPattern,
				},
			},
		},
		MaxItems:   100, // May need adjustment
		StartToken: "",  // For pagination, start at first page
	}
	req := connect.NewRequest(queryReq)
	return u.SearchUtxosWithContext(ctx, req)
}

func (u *UtxorpcClient) SubmitTransaction(
	txCbor string,
) (*connect.Response[submit.SubmitTxResponse], error) {
	ctx := context.Background()
	// Decode the transaction data from hex
	txRawBytes, err := hex.DecodeString(txCbor)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction hash: %w", err)
	}

	// Create a SubmitTxRequest with the transaction data
	tx := &submit.AnyChainTx{
		Type: &submit.AnyChainTx_Raw{
			Raw: txRawBytes,
		},
	}

	// Create a list with one transaction
	req := &submit.SubmitTxRequest{
		Tx: []*submit.AnyChainTx{tx},
	}
	return u.SubmitTransactionWithContext(ctx, req)
}

func (u *UtxorpcClient) SubmitTransactionWithContext(
	ctx context.Context,
	txReq *submit.SubmitTxRequest,
) (*connect.Response[submit.SubmitTxResponse], error) {
	req := connect.NewRequest(txReq)
	u.AddHeadersToRequest(req)
	return u.Submit.SubmitTx(ctx, req)
}

func (u *UtxorpcClient) WaitForTransaction(
	txRef string,
) (*connect.ServerStreamForClient[submit.WaitForTxResponse], error) {
	ctx := context.Background()
	// Decode the transaction references from hex
	var decodedRefs [][]byte
	refBytes, err := hex.DecodeString(txRef)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode transaction reference %s: %w",
			txRef,
			err,
		)
	}
	decodedRefs = append(decodedRefs, refBytes)

	// Create a WaitForTxRequest with the decoded transaction references
	req := &submit.WaitForTxRequest{
		Ref: decodedRefs,
	}
	return u.WaitForTransactionWithContext(ctx, req)
}

func (u *UtxorpcClient) WaitForTransactionWithContext(
	ctx context.Context,
	txReq *submit.WaitForTxRequest,
) (*connect.ServerStreamForClient[submit.WaitForTxResponse], error) {
	req := connect.NewRequest(txReq)
	u.AddHeadersToRequest(req)
	return u.Submit.WaitForTx(ctx, req)
}

func (u *UtxorpcClient) WatchMempoolTransactions() (
	*connect.ServerStreamForClient[submit.WatchMempoolResponse],
	error,
) {
	ctx := context.Background()
	return u.WatchMempoolTransactionsWithContext(ctx)
}

func (u *UtxorpcClient) WatchMempoolTransactionsWithContext(ctx context.Context) (
	*connect.ServerStreamForClient[submit.WatchMempoolResponse],
	error,
) {
	req := connect.NewRequest(&submit.WatchMempoolRequest{})
	u.AddHeadersToRequest(req)
	return u.Submit.WatchMempool(ctx, req)
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

func (u *UtxorpcClient) GetBlockByRef(
	blockHashStr string,
	blockIndex int64,
) (*connect.Response[sync.FetchBlockResponse], error) {
	ctx := context.Background()
	req := &sync.FetchBlockRequest{Ref: syncIntersect(blockHashStr, blockIndex)}
	return u.GetBlockByRefWithContext(ctx, req)
}

func (u *UtxorpcClient) GetBlockByRefWithContext(
	ctx context.Context,
	blockReq *sync.FetchBlockRequest,
) (*connect.Response[sync.FetchBlockResponse], error) {
	req := connect.NewRequest(blockReq)
	u.AddHeadersToRequest(req)
	return u.Sync.FetchBlock(ctx, req)
}

func (u *UtxorpcClient) WatchBlocksByRef(
	blockHashStr string,
	blockIndex int64,
) (*connect.ServerStreamForClient[sync.FollowTipResponse], error) {
	ctx := context.Background()
	req := &sync.FollowTipRequest{
		Intersect: syncIntersect(blockHashStr, blockIndex),
	}
	return u.WatchBlocksByRefWithContext(ctx, req)
}

func (u *UtxorpcClient) WatchBlocksByRefWithContext(
	ctx context.Context,
	blockReq *sync.FollowTipRequest,
) (*connect.ServerStreamForClient[sync.FollowTipResponse], error) {
	req := connect.NewRequest(blockReq)
	u.AddHeadersToRequest(req)
	return u.Sync.FollowTip(ctx, req)
}

func (u *UtxorpcClient) GetTip() (*connect.Response[sync.ReadTipResponse], error) {
	return u.GetTipWithContext(context.Background())
}

func (u *UtxorpcClient) GetTipWithContext(
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
	fetchBlockReqProto := &sync.FetchBlockRequest{
		Ref: []*sync.BlockRef{blockRef},
	}
	reqFetchBlock := connect.NewRequest(fetchBlockReqProto)
	u.AddHeadersToRequest(reqFetchBlock)

	blockRespFull, err := u.Sync.FetchBlock(ctx, reqFetchBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block for tip: %w", err)
	}
	if blockRespFull.Msg == nil || len(blockRespFull.Msg.GetBlock()) == 0 ||
		blockRespFull.Msg.GetBlock()[0] == nil {
		return nil, errors.New(
			"received nil or empty block data from FetchBlockResponse for tip",
		)
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

func (u *UtxorpcClient) WatchTransaction(
	blockHashStr string,
	blockIndex int64,
) (*connect.ServerStreamForClient[watch.WatchTxResponse], error) {
	ctx := context.Background()
	req := &watch.WatchTxRequest{
		Intersect: watchIntersect(blockHashStr, blockIndex),
	}
	return u.WatchTransactionWithContext(ctx, req)
}

func (u *UtxorpcClient) WatchTransactionWithContext(
	ctx context.Context,
	watchReq *watch.WatchTxRequest,
) (*connect.ServerStreamForClient[watch.WatchTxResponse], error) {
	req := connect.NewRequest(watchReq)
	u.AddHeadersToRequest(req)
	return u.Watch.WatchTx(ctx, req)
}
