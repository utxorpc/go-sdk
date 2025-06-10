package sdk

import (
	"context"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/submit"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/submit/submitconnect"
)

type SubmitServiceClient = submitconnect.SubmitServiceClient

func (u *UtxorpcClient) NewSubmitServiceClient() SubmitServiceClient {
	return submitconnect.NewSubmitServiceClient(
		u.httpClient,
		u.baseUrl,
		connect.WithGRPC(),
	)
}

func (u *UtxorpcClient) EvalTx(
	req *connect.Request[submit.EvalTxRequest],
) (*connect.Response[submit.EvalTxResponse], error) {
	ctx := context.Background()
	return u.EvalTxWithContext(ctx, req)
}

func (u *UtxorpcClient) EvalTxWithContext(
	ctx context.Context,
	req *connect.Request[submit.EvalTxRequest],
) (*connect.Response[submit.EvalTxResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Submit.EvalTx(ctx, req)
}

func (u *UtxorpcClient) ReadMempool(
	req *connect.Request[submit.ReadMempoolRequest],
) (*connect.Response[submit.ReadMempoolResponse], error) {
	ctx := context.Background()
	return u.ReadMempoolWithContext(ctx, req)
}

func (u *UtxorpcClient) ReadMempoolWithContext(
	ctx context.Context,
	req *connect.Request[submit.ReadMempoolRequest],
) (*connect.Response[submit.ReadMempoolResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Submit.ReadMempool(ctx, req)
}

func (u *UtxorpcClient) SubmitTx(
	req *connect.Request[submit.SubmitTxRequest],
) (*connect.Response[submit.SubmitTxResponse], error) {
	ctx := context.Background()
	return u.SubmitTxWithContext(ctx, req)
}

func (u *UtxorpcClient) SubmitTxWithContext(
	ctx context.Context,
	req *connect.Request[submit.SubmitTxRequest],
) (*connect.Response[submit.SubmitTxResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Submit.SubmitTx(ctx, req)
}

func (u *UtxorpcClient) WaitForTx(
	req *connect.Request[submit.WaitForTxRequest],
) (*connect.ServerStreamForClient[submit.WaitForTxResponse], error) {
	ctx := context.Background()
	return u.WaitForTxWithContext(ctx, req)
}

func (u *UtxorpcClient) WaitForTxWithContext(
	ctx context.Context,
	req *connect.Request[submit.WaitForTxRequest],
) (*connect.ServerStreamForClient[submit.WaitForTxResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Submit.WaitForTx(ctx, req)
}

func (u *UtxorpcClient) WatchMempool(
	req *connect.Request[submit.WatchMempoolRequest],
) (*connect.ServerStreamForClient[submit.WatchMempoolResponse], error) {
	ctx := context.Background()
	return u.WatchMempoolWithContext(ctx, req)
}

func (u *UtxorpcClient) WatchMempoolWithContext(
	ctx context.Context,
	req *connect.Request[submit.WatchMempoolRequest],
) (*connect.ServerStreamForClient[submit.WatchMempoolResponse], error) {
	u.AddHeadersToRequest(req)
	return u.Submit.WatchMempool(ctx, req)
}
