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

func (u *UtxorpcClient) ReadData(
	req *query.ReadDataRequest,
) (*connect.Response[query.ReadDataResponse], error) {
	ctx := context.Background()
	return u.ReadDataWithContext(ctx, req)
}

func (u *UtxorpcClient) ReadDataWithContext(
	ctx context.Context,
	queryReq *query.ReadDataRequest,
) (*connect.Response[query.ReadDataResponse], error) {
	req := connect.NewRequest(queryReq)
	u.AddHeadersToRequest(req)
	return u.Query.ReadData(ctx, req)
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
	req *query.SearchUtxosRequest,
) (*connect.Response[query.SearchUtxosResponse], error) {
	ctx := context.Background()
	return u.SearchUtxosWithContext(ctx, req)
}

func (u *UtxorpcClient) SearchUtxosWithContext(
	ctx context.Context,
	queryReq *query.SearchUtxosRequest,
) (*connect.Response[query.SearchUtxosResponse], error) {
	req := connect.NewRequest(queryReq)
	u.AddHeadersToRequest(req)
	return u.Query.SearchUtxos(ctx, req)
}

// Helpers

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
	req := &query.ReadUtxosRequest{Keys: []*query.TxoRef{txoRef}}
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

	req := &query.ReadUtxosRequest{Keys: refs}
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

	return u.SearchUtxosWithContext(ctx, queryReq)
}

func (u *UtxorpcClient) GetUtxosByAddressWithAsset(
	addressBytes []byte,
	policyIdBytes []byte,
	assetNameBytes []byte,
) (*connect.Response[query.SearchUtxosResponse], error) {
	return u.GetUtxosByAddressWithAssetWithContext(context.Background(), addressBytes, policyIdBytes, assetNameBytes)
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

	return u.SearchUtxosWithContext(ctx, queryReq)
}

func (u *UtxorpcClient) GetUtxosByAsset(
	policyIdBytes []byte,
	assetNameBytes []byte,
) (*connect.Response[query.SearchUtxosResponse], error) {
	return u.GetUtxosByAssetWithContext(context.Background(), policyIdBytes, assetNameBytes)
}

func (u *UtxorpcClient) GetUtxosByAssetWithContext(
	ctx context.Context,
	policyIdBytes []byte,
	assetNameBytes []byte,
) (*connect.Response[query.SearchUtxosResponse], error) {
	if policyIdBytes == nil && assetNameBytes == nil {
		return nil, errors.New("at least one of policyId or assetName must be provided")
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

	return u.SearchUtxosWithContext(ctx, queryReq)
}
