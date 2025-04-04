package sdk

import (
	"context"
	"encoding/base64"
	"encoding/hex"

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

func (u *UtxorpcClient) ReadUtxo(
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
	return u.ReadUtxos(req)
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

func (u *UtxorpcClient) GetUtxosByAddress(
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
	return u.SearchUtxos(queryReq)
}
