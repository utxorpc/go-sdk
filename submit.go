package sdk

import (
	"context"
	"encoding/hex"
	"fmt"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/submit"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/submit/submitconnect"
)

type SubmitServiceClient submitconnect.SubmitServiceClient

func (u *UtxorpcClient) NewSubmitServiceClient() SubmitServiceClient {
	return submitconnect.NewSubmitServiceClient(
		u.httpClient,
		u.baseUrl,
		connect.WithGRPC(),
	)
}

func (u *UtxorpcClient) EvalTx(
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
	return u.EvalTxWithContext(ctx, req)
}

func (u *UtxorpcClient) EvalTxWithContext(
	ctx context.Context,
	txReq *submit.EvalTxRequest,
) (*connect.Response[submit.EvalTxResponse], error) {
	req := connect.NewRequest(txReq)
	u.AddHeadersToRequest(req)
	return u.Submit.EvalTx(ctx, req)
}

func (u *UtxorpcClient) ReadMempool() (*connect.Response[submit.ReadMempoolResponse], error) {
	ctx := context.Background()
	return u.ReadMempoolWithContext(ctx)
}

func (u *UtxorpcClient) ReadMempoolWithContext(
	ctx context.Context,
) (*connect.Response[submit.ReadMempoolResponse], error) {
	req := connect.NewRequest(&submit.ReadMempoolRequest{})
	u.AddHeadersToRequest(req)
	return u.Submit.ReadMempool(ctx, req)
}

func (u *UtxorpcClient) SubmitTx(
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
	return u.SubmitTxWithContext(ctx, req)
}

func (u *UtxorpcClient) SubmitTxWithContext(
	ctx context.Context,
	txReq *submit.SubmitTxRequest,
) (*connect.Response[submit.SubmitTxResponse], error) {
	req := connect.NewRequest(txReq)
	u.AddHeadersToRequest(req)
	return u.Submit.SubmitTx(ctx, req)
}

func (u *UtxorpcClient) WaitForTx(
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
	return u.WaitForTxWithContext(ctx, req)
}

func (u *UtxorpcClient) WaitForTxWithContext(
	ctx context.Context,
	txReq *submit.WaitForTxRequest,
) (*connect.ServerStreamForClient[submit.WaitForTxResponse], error) {
	req := connect.NewRequest(txReq)
	u.AddHeadersToRequest(req)
	return u.Submit.WaitForTx(ctx, req)
}

func (u *UtxorpcClient) WatchMempool() (
	*connect.ServerStreamForClient[submit.WatchMempoolResponse],
	error,
) {
	ctx := context.Background()
	return u.WatchMempoolWithContext(ctx)
}

func (u *UtxorpcClient) WatchMempoolWithContext(ctx context.Context) (
	*connect.ServerStreamForClient[submit.WatchMempoolResponse],
	error,
) {
	req := connect.NewRequest(&submit.WatchMempoolRequest{})
	u.AddHeadersToRequest(req)
	return u.Submit.WatchMempool(ctx, req)
}
