package sdk

import (
	"context"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/query"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/query/queryconnect"
)

type QueryServiceClient = queryconnect.QueryServiceClient

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
	req *connect.Request[query.ReadDataRequest],
) (*connect.Response[query.ReadDataResponse], error) {
	ctx := context.Background()
	return u.ReadDataWithContext(ctx, req)
}

func (u *UtxorpcClient) ReadDataWithContext(
	ctx context.Context,
	req *connect.Request[query.ReadDataRequest],
) (*connect.Response[query.ReadDataResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Query.ReadData(ctx, req)
}

func (u *UtxorpcClient) ReadParams(
	req *connect.Request[query.ReadParamsRequest],
) (*connect.Response[query.ReadParamsResponse], error) {
	ctx := context.Background()
	return u.ReadParamsWithContext(ctx, req)
}

func (u *UtxorpcClient) ReadParamsWithContext(
	ctx context.Context,
	req *connect.Request[query.ReadParamsRequest],
) (*connect.Response[query.ReadParamsResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Query.ReadParams(ctx, req)
}

func (u *UtxorpcClient) ReadUtxos(
	req *connect.Request[query.ReadUtxosRequest],
) (*connect.Response[query.ReadUtxosResponse], error) {
	ctx := context.Background()
	return u.ReadUtxosWithContext(ctx, req)
}

func (u *UtxorpcClient) ReadUtxosWithContext(
	ctx context.Context,
	req *connect.Request[query.ReadUtxosRequest],
) (*connect.Response[query.ReadUtxosResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Query.ReadUtxos(ctx, req)
}

func (u *UtxorpcClient) SearchUtxos(
	req *connect.Request[query.SearchUtxosRequest],
) (*connect.Response[query.SearchUtxosResponse], error) {
	ctx := context.Background()
	return u.SearchUtxosWithContext(ctx, req)
}

func (u *UtxorpcClient) SearchUtxosWithContext(
	ctx context.Context,
	req *connect.Request[query.SearchUtxosRequest],
) (*connect.Response[query.SearchUtxosResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Query.SearchUtxos(ctx, req)
}
