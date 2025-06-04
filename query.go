package sdk

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/cardano"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/query"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/query/queryconnect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type QueryServiceClient queryconnect.QueryServiceClient

func NewQueryServiceClient(u *UtxorpcClient) QueryServiceClient {
	return u.NewQueryServiceClient()
}

func (u *UtxorpcClient) NewQueryServiceClient() QueryServiceClient {
	return queryconnect.NewQueryServiceClient(
		u.httpClient,
		u.baseUrl,
		connect.WithGRPC(),
	)
}

func (u *UtxorpcClient) QueryService() QueryServiceClient {
	return u.Query
}

func (u *UtxorpcClient) ReadParams() (*connect.Response[query.ReadParamsResponse], error) {
	ctx := context.Background()
	return u.ReadParamsWithContext(ctx)
}

func (u *UtxorpcClient) ReadParamsWithContext(
	ctx context.Context,
) (*connect.Response[query.ReadParamsResponse], error) {
	req := connect.NewRequest(&query.ReadParamsRequest{})
	u.AddHeadersToRequest(req)
	return u.Query.ReadParams(ctx, req)
}

func (u *UtxorpcClient) ReadUtxos(
	req *query.ReadUtxosRequest,
) (*connect.Response[query.ReadUtxosResponse], error) {
	ctx := context.Background()
	return u.ReadUtxosWithContext(ctx, req)
}

func (u *UtxorpcClient) ReadUtxosWithContext(
	ctx context.Context,
	queryReq *query.ReadUtxosRequest,
) (*connect.Response[query.ReadUtxosResponse], error) {
	req := connect.NewRequest(queryReq)
	u.AddHeadersToRequest(req)
	return u.Query.ReadUtxos(ctx, req)
}

func (u *UtxorpcClient) SearchUtxos(
	predicate *query.UtxoPredicate,
	maxItems int32,
	startToken string,
	fieldMask *fieldmaskpb.FieldMask,
) (*connect.Response[query.SearchUtxosResponse], error) {
	ctx := context.Background()
	queryReq := &query.SearchUtxosRequest{
		Predicate:  predicate,
		MaxItems:   maxItems,
		StartToken: startToken,
		FieldMask:  fieldMask,
	}
	req := connect.NewRequest(queryReq)
	u.AddHeadersToRequest(req)
	return u.Query.SearchUtxos(ctx, req)
}

func (u *UtxorpcClient) SearchUtxosWithContext(
	ctx context.Context,
	predicate *query.UtxoPredicate,
	maxItems int32,
	startToken string,
	fieldMask *fieldmaskpb.FieldMask,
) (*connect.Response[query.SearchUtxosResponse], error) {
	queryReq := &query.SearchUtxosRequest{
		Predicate:  predicate,
		MaxItems:   maxItems,
		StartToken: startToken,
		FieldMask:  fieldMask,
	}
	req := connect.NewRequest(queryReq)
	u.AddHeadersToRequest(req)
	return u.Query.SearchUtxos(ctx, req)
}

// Helpers

func (u *UtxorpcClient) GetUtxoByRef(
	txHashStr string,
	txIndex uint32,
	ctx ...context.Context,
) (*connect.Response[query.ReadUtxosResponse], error) {
	var ctxToUse context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		ctxToUse = ctx[0]
	} else {
		ctxToUse = context.Background()
	}

	var txHashBytes []byte
	var err error
	txHashBytes, hexErr := hex.DecodeString(txHashStr)
	if hexErr != nil {
		txHashBytes, err = base64.StdEncoding.DecodeString(txHashStr)
		if err != nil {
			return nil, err
		}
	}
	txoRef := &query.TxoRef{
		Hash:  txHashBytes,
		Index: txIndex,
	}
	req := &query.ReadUtxosRequest{Keys: []*query.TxoRef{txoRef}}
	return u.ReadUtxosWithContext(ctxToUse, req)
}

type TxoReference struct {
	TxHash string
	Index  uint32
}

func (u *UtxorpcClient) GetUtxosByRefs(
	refs []TxoReference,
	batchSize *int,
	ctx ...context.Context,
) (*connect.Response[query.ReadUtxosResponse], error) {
	var ctxToUse context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		ctxToUse = ctx[0]
	} else {
		ctxToUse = context.Background()
	}

	if len(refs) == 0 {
		return nil, errors.New("no transaction references provided")
	}

	defaultBatchSize := 100
	if batchSize != nil && *batchSize > 0 {
		defaultBatchSize = *batchSize
	}

	allTxoRefs := make([]*query.TxoRef, 0, len(refs))
	for _, ref := range refs {
		var txHashBytes []byte
		txHashBytes, err := hex.DecodeString(ref.TxHash)
		if err != nil {
			return nil, err
		}
		txoRef := &query.TxoRef{
			Hash:  txHashBytes,
			Index: ref.Index,
		}
		allTxoRefs = append(allTxoRefs, txoRef)
	}

	if len(allTxoRefs) <= defaultBatchSize {
		req := &query.ReadUtxosRequest{Keys: allTxoRefs}
		return u.ReadUtxosWithContext(ctxToUse, req)
	}

	var allResults []*query.AnyUtxoData
	for i := 0; i < len(allTxoRefs); i += defaultBatchSize {
		end := i + defaultBatchSize
		if end > len(allTxoRefs) {
			end = len(allTxoRefs)
		}

		batch := allTxoRefs[i:end]
		req := &query.ReadUtxosRequest{Keys: batch}

		resp, err := u.ReadUtxosWithContext(ctxToUse, req)
		if err != nil {
			return nil, err
		}

		if resp.Msg != nil && resp.Msg.Items != nil {
			allResults = append(allResults, resp.Msg.Items...)
		}
	}

	aggregatedResponse := &query.ReadUtxosResponse{
		Items: allResults,
	}

	return connect.NewResponse(aggregatedResponse), nil
}

func (u *UtxorpcClient) GetUtxosByAddress(
	address []byte,
	maxItems *int32,
	startToken *string,
	ctx ...context.Context,
) (*connect.Response[query.SearchUtxosResponse], error) {
	var ctxToUse context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		ctxToUse = ctx[0]
	} else {
		ctxToUse = context.Background()
	}

	items := int32(100)
	if maxItems != nil {
		items = *maxItems
	}

	token := ""
	if startToken != nil {
		token = *startToken
	}

	predicate := &query.UtxoPredicate{
		Match: &query.AnyUtxoPattern{
			UtxoPattern: &query.AnyUtxoPattern_Cardano{
				Cardano: &cardano.TxOutputPattern{
					Address: &cardano.AddressPattern{
						ExactAddress: address,
					},
				},
			},
		},
	}
	fieldMask := &fieldmaskpb.FieldMask{Paths: []string{}}
	return u.SearchUtxosWithContext(ctxToUse, predicate, items, token, fieldMask)
}

func (u *UtxorpcClient) GetUtxosByAddressWithAsset(
	addressBytes []byte,
	policyIdBytes []byte,
	assetNameBytes []byte,
	maxItems *int32,
	startToken *string,
	ctx ...context.Context,
) (*connect.Response[query.SearchUtxosResponse], error) {
	var ctxToUse context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		ctxToUse = ctx[0]
	} else {
		ctxToUse = context.Background()
	}

	items := int32(100)
	if maxItems != nil && *maxItems > 0 {
		items = *maxItems
	}
	token := ""
	if startToken != nil {
		token = *startToken
	}

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

	predicate := &query.UtxoPredicate{
		Match: &query.AnyUtxoPattern{
			UtxoPattern: &query.AnyUtxoPattern_Cardano{
				Cardano: tpl,
			},
		},
	}

	return u.SearchUtxosWithContext(ctxToUse, predicate, items, token, nil)
}

func (u *UtxorpcClient) GetUtxosByAsset(
	policyIdBytes []byte,
	assetNameBytes []byte,
	maxItems *int32,
	startToken *string,
	ctx ...context.Context,
) (*connect.Response[query.SearchUtxosResponse], error) {
	var ctxToUse context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		ctxToUse = ctx[0]
	} else {
		ctxToUse = context.Background()
	}

	if policyIdBytes == nil && assetNameBytes == nil {
		return nil, errors.New("at least one of policyId or assetName must be provided")
	}

	items := int32(100)
	if maxItems != nil {
		items = *maxItems
	}

	token := ""
	if startToken != nil {
		token = *startToken
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

	predicate := &query.UtxoPredicate{
		Match: &query.AnyUtxoPattern{
			UtxoPattern: &query.AnyUtxoPattern_Cardano{
				Cardano: cardanoOutputPattern,
			},
		},
	}
	return u.SearchUtxosWithContext(ctxToUse, predicate, items, token, nil)
}
