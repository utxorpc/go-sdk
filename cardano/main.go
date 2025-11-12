package cardano

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
	sdk "github.com/utxorpc/go-sdk"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type Client struct {
	UtxorpcClient *sdk.UtxorpcClient
}

func NewClient(options ...sdk.ClientOption) *Client {
	c := &Client{}
	c.UtxorpcClient = sdk.NewClient(options...)
	return c
}

func (c *Client) GetProtocolParameters() (*connect.Response[query.ReadParamsResponse], error) {
	ctx := context.Background()
	return c.GetProtocolParametersWithContext(ctx)
}

func (c *Client) GetProtocolParametersWithContext(
	ctx context.Context,
) (*connect.Response[query.ReadParamsResponse], error) {
	req := connect.NewRequest(&query.ReadParamsRequest{})
	c.UtxorpcClient.AddHeadersToRequest(req)
	return c.UtxorpcClient.Query.ReadParams(ctx, req)
}

func (c *Client) GetUtxoByRef(
	txHashStr string,
	txIndex uint32,
) (*connect.Response[query.ReadUtxosResponse], error) {
	return c.GetUtxoByRefWithContext(context.Background(), txHashStr, txIndex)
}

func (c *Client) GetUtxoByRefWithContext(
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
	return c.UtxorpcClient.ReadUtxosWithContext(ctx, req)
}

func (c *Client) EvaluateTransaction(
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
		Tx: tx,
	}
	return c.EvaluateTransactionWithContext(ctx, req)
}

func (c *Client) EvaluateTransactionWithContext(
	ctx context.Context,
	txReq *submit.EvalTxRequest,
) (*connect.Response[submit.EvalTxResponse], error) {
	req := connect.NewRequest(txReq)
	c.UtxorpcClient.AddHeadersToRequest(req)
	return c.UtxorpcClient.Submit.EvalTx(ctx, req)
}

func (c *Client) GetMempoolTransactions() (*connect.Response[submit.ReadMempoolResponse], error) {
	ctx := context.Background()
	return c.GetMempoolTransactionsWithContext(ctx)
}

func (c *Client) GetMempoolTransactionsWithContext(
	ctx context.Context,
) (*connect.Response[submit.ReadMempoolResponse], error) {
	req := connect.NewRequest(&submit.ReadMempoolRequest{})
	c.UtxorpcClient.AddHeadersToRequest(req)
	return c.UtxorpcClient.Submit.ReadMempool(ctx, req)
}

func (c *Client) GetUtxosByRefs(
	refs []*query.TxoRef,
) (*connect.Response[query.ReadUtxosResponse], error) {
	return c.GetUtxosByRefsWithContext(context.Background(), refs)
}

func (c *Client) GetUtxosByRefsWithContext(
	ctx context.Context,
	refs []*query.TxoRef,
) (*connect.Response[query.ReadUtxosResponse], error) {
	if len(refs) == 0 {
		return nil, errors.New("no transaction references provided")
	}

	txReq := &query.ReadUtxosRequest{Keys: refs}
	req := connect.NewRequest(txReq)
	return c.UtxorpcClient.ReadUtxosWithContext(ctx, req)
}

func (c *Client) GetUtxosByAddress(
	address []byte,
) (*connect.Response[query.SearchUtxosResponse], error) {
	return c.GetUtxosByAddressWithContext(context.Background(), address)
}

func (c *Client) GetUtxosByAddressWithContext(
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
	return c.UtxorpcClient.SearchUtxosWithContext(ctx, req)
}

func (c *Client) GetUtxosByAddressWithAsset(
	addressBytes []byte,
	policyIdBytes []byte,
	assetNameBytes []byte,
) (*connect.Response[query.SearchUtxosResponse], error) {
	return c.GetUtxosByAddressWithAssetWithContext(
		context.Background(),
		addressBytes,
		policyIdBytes,
		assetNameBytes,
	)
}

func (c *Client) GetUtxosByAddressWithAssetWithContext(
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
	return c.UtxorpcClient.SearchUtxosWithContext(ctx, req)
}

func (c *Client) GetUtxosByAsset(
	policyIdBytes []byte,
	assetNameBytes []byte,
) (*connect.Response[query.SearchUtxosResponse], error) {
	return c.GetUtxosByAssetWithContext(
		context.Background(),
		policyIdBytes,
		assetNameBytes,
	)
}

func (c *Client) GetUtxosByAssetWithContext(
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
	return c.UtxorpcClient.SearchUtxosWithContext(ctx, req)
}

func (c *Client) SubmitTransaction(
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
		Tx: tx,
	}
	return c.SubmitTransactionWithContext(ctx, req)
}

func (c *Client) SubmitTransactionWithContext(
	ctx context.Context,
	txReq *submit.SubmitTxRequest,
) (*connect.Response[submit.SubmitTxResponse], error) {
	req := connect.NewRequest(txReq)
	c.UtxorpcClient.AddHeadersToRequest(req)
	return c.UtxorpcClient.Submit.SubmitTx(ctx, req)
}

func (c *Client) WaitForTransaction(
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
	return c.WaitForTransactionWithContext(ctx, req)
}

func (c *Client) WaitForTransactionWithContext(
	ctx context.Context,
	txReq *submit.WaitForTxRequest,
) (*connect.ServerStreamForClient[submit.WaitForTxResponse], error) {
	req := connect.NewRequest(txReq)
	c.UtxorpcClient.AddHeadersToRequest(req)
	return c.UtxorpcClient.Submit.WaitForTx(ctx, req)
}

func (c *Client) WatchMempoolTransactions() (
	*connect.ServerStreamForClient[submit.WatchMempoolResponse],
	error,
) {
	ctx := context.Background()
	return c.WatchMempoolTransactionsWithContext(ctx)
}

func (c *Client) WatchMempoolTransactionsWithContext(ctx context.Context) (
	*connect.ServerStreamForClient[submit.WatchMempoolResponse],
	error,
) {
	req := connect.NewRequest(&submit.WatchMempoolRequest{})
	c.UtxorpcClient.AddHeadersToRequest(req)
	return c.UtxorpcClient.Submit.WatchMempool(ctx, req)
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
		blockRef.Slot = uint64(blockIndex)
	}
	// Only add blockRef to intersect if at least one of blockHashStr or blockIndex is provided
	if blockHashStr != "" || blockIndex > -1 {
		intersect = []*sync.BlockRef{blockRef}
	}
	return intersect
}

func (c *Client) GetBlockByRef(
	blockHashStr string,
	blockIndex int64,
) (*connect.Response[sync.FetchBlockResponse], error) {
	ctx := context.Background()
	req := &sync.FetchBlockRequest{Ref: syncIntersect(blockHashStr, blockIndex)}
	return c.GetBlockByRefWithContext(ctx, req)
}

func (c *Client) GetBlockByRefWithContext(
	ctx context.Context,
	blockReq *sync.FetchBlockRequest,
) (*connect.Response[sync.FetchBlockResponse], error) {
	req := connect.NewRequest(blockReq)
	c.UtxorpcClient.AddHeadersToRequest(req)
	return c.UtxorpcClient.Sync.FetchBlock(ctx, req)
}

func (c *Client) WatchBlocksByRef(
	blockHashStr string,
	blockIndex int64,
) (*connect.ServerStreamForClient[sync.FollowTipResponse], error) {
	ctx := context.Background()
	req := &sync.FollowTipRequest{
		Intersect: syncIntersect(blockHashStr, blockIndex),
	}
	return c.WatchBlocksByRefWithContext(ctx, req)
}

func (c *Client) WatchBlocksByRefWithContext(
	ctx context.Context,
	blockReq *sync.FollowTipRequest,
) (*connect.ServerStreamForClient[sync.FollowTipResponse], error) {
	req := connect.NewRequest(blockReq)
	c.UtxorpcClient.AddHeadersToRequest(req)
	return c.UtxorpcClient.Sync.FollowTip(ctx, req)
}

func (c *Client) GetTip() (*connect.Response[sync.ReadTipResponse], error) {
	return c.GetTipWithContext(context.Background())
}

func (c *Client) GetTipWithContext(
	ctx context.Context,
) (*connect.Response[sync.ReadTipResponse], error) {
	readTipReqProto := &sync.ReadTipRequest{}
	reqReadTip := connect.NewRequest(readTipReqProto)
	c.UtxorpcClient.AddHeadersToRequest(reqReadTip)

	tipResp, err := c.UtxorpcClient.Sync.ReadTip(ctx, reqReadTip)
	if err != nil {
		return nil, fmt.Errorf("failed to read tip: %w", err)
	}
	if tipResp.Msg == nil || tipResp.Msg.GetTip() == nil {
		return nil, errors.New("received nil tip from ReadTipResponse")
	}

	return tipResp, nil
}

func (c *Client) ReadBlock(
	blockRef *sync.BlockRef,
) (*connect.Response[sync.FetchBlockResponse], error) {
	return c.ReadBlockWithContext(context.Background(), blockRef)
}

func (c *Client) ReadBlockWithContext(
	ctx context.Context,
	blockRef *sync.BlockRef,
) (*connect.Response[sync.FetchBlockResponse], error) {
	fetchBlockReqProto := &sync.FetchBlockRequest{
		Ref: []*sync.BlockRef{blockRef},
	}
	reqFetchBlock := connect.NewRequest(fetchBlockReqProto)
	c.UtxorpcClient.AddHeadersToRequest(reqFetchBlock)

	blockRespFull, err := c.UtxorpcClient.Sync.FetchBlock(ctx, reqFetchBlock)
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
		blockRef.Slot = uint64(blockIndex)
	}
	// Only add blockRef to intersect if at least one of blockHashStr or blockIndex is provided
	if blockHashStr != "" || blockIndex > -1 {
		intersect = []*watch.BlockRef{blockRef}
	}
	return intersect
}

func (c *Client) WatchTransaction(
	blockHashStr string,
	blockIndex int64,
) (*connect.ServerStreamForClient[watch.WatchTxResponse], error) {
	ctx := context.Background()
	req := &watch.WatchTxRequest{
		Intersect: watchIntersect(blockHashStr, blockIndex),
	}
	return c.WatchTransactionWithContext(ctx, req)
}

func (c *Client) WatchTransactionWithContext(
	ctx context.Context,
	watchReq *watch.WatchTxRequest,
) (*connect.ServerStreamForClient[watch.WatchTxResponse], error) {
	req := connect.NewRequest(watchReq)
	c.UtxorpcClient.AddHeadersToRequest(req)
	return c.UtxorpcClient.Watch.WatchTx(ctx, req)
}
