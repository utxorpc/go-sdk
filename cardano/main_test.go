package cardano

import (
	"bytes"
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1beta/query"
	sdk "github.com/utxorpc/go-sdk"
)

func TestGetUtxoByRefBuildsReadUtxosRequest(t *testing.T) {
	fakeQuery := &recordingQueryClient{}
	client := NewClient(
		sdk.WithBaseUrl("http://example.test"),
		sdk.WithHeaders(map[string]string{"dmtr-api-key": "secret"}),
	)
	client.UtxorpcClient.Query = fakeQuery

	_, err := client.GetUtxoByRefWithContext(
		context.Background(),
		"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f",
		7,
	)
	if err != nil {
		t.Fatalf("GetUtxoByRefWithContext returned error: %v", err)
	}
	if fakeQuery.readUtxosReq == nil {
		t.Fatal("ReadUtxos was not called")
	}
	if got := fakeQuery.readUtxosReq.Header().Get("dmtr-api-key"); got != "secret" {
		t.Fatalf("dmtr-api-key header = %q, want %q", got, "secret")
	}

	keys := fakeQuery.readUtxosReq.Msg.GetKeys()
	if len(keys) != 1 {
		t.Fatalf("len(keys) = %d, want 1", len(keys))
	}
	if got := keys[0].GetIndex(); got != 7 {
		t.Fatalf("txo index = %d, want 7", got)
	}
	wantHash := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}
	if got := keys[0].GetHash(); !bytes.Equal(got, wantHash) {
		t.Fatalf("txo hash = %x, want %x", got, wantHash)
	}
}

func TestGetUtxosByAssetBuildsSearchRequest(t *testing.T) {
	fakeQuery := &recordingQueryClient{}
	client := NewClient(sdk.WithBaseUrl("http://example.test"))
	client.UtxorpcClient.Query = fakeQuery

	policyID := []byte{0x01, 0x02, 0x03}
	assetName := []byte("asset")

	_, err := client.GetUtxosByAssetWithContext(
		context.Background(),
		policyID,
		assetName,
	)
	if err != nil {
		t.Fatalf("GetUtxosByAssetWithContext returned error: %v", err)
	}
	if fakeQuery.searchUtxosReq == nil {
		t.Fatal("SearchUtxos was not called")
	}

	req := fakeQuery.searchUtxosReq.Msg
	if got := req.GetMaxItems(); got != 100 {
		t.Fatalf("max_items = %d, want 100", got)
	}
	if got := req.GetStartToken(); got != "" {
		t.Fatalf("start_token = %q, want empty", got)
	}
	if req.GetFieldMask() == nil {
		t.Fatal("field_mask is nil")
	}

	pattern := req.GetPredicate().GetMatch().GetCardano()
	if pattern == nil {
		t.Fatal("cardano UTxO pattern is nil")
	}
	asset := pattern.GetAsset()
	if asset == nil {
		t.Fatal("asset pattern is nil")
	}
	if got := asset.GetPolicyId(); !bytes.Equal(got, policyID) {
		t.Fatalf("policy_id = %x, want %x", got, policyID)
	}
	if got := asset.GetAssetName(); !bytes.Equal(got, assetName) {
		t.Fatalf("asset_name = %x, want %x", got, assetName)
	}
}

type recordingQueryClient struct {
	readUtxosReq   *connect.Request[query.ReadUtxosRequest]
	searchUtxosReq *connect.Request[query.SearchUtxosRequest]
}

func (*recordingQueryClient) ReadParams(
	context.Context,
	*connect.Request[query.ReadParamsRequest],
) (*connect.Response[query.ReadParamsResponse], error) {
	return connect.NewResponse(&query.ReadParamsResponse{}), nil
}

func (r *recordingQueryClient) ReadUtxos(
	_ context.Context,
	req *connect.Request[query.ReadUtxosRequest],
) (*connect.Response[query.ReadUtxosResponse], error) {
	r.readUtxosReq = req
	return connect.NewResponse(&query.ReadUtxosResponse{}), nil
}

func (r *recordingQueryClient) SearchUtxos(
	_ context.Context,
	req *connect.Request[query.SearchUtxosRequest],
) (*connect.Response[query.SearchUtxosResponse], error) {
	r.searchUtxosReq = req
	return connect.NewResponse(&query.SearchUtxosResponse{}), nil
}

func (*recordingQueryClient) ReadData(
	context.Context,
	*connect.Request[query.ReadDataRequest],
) (*connect.Response[query.ReadDataResponse], error) {
	return connect.NewResponse(&query.ReadDataResponse{}), nil
}

func (*recordingQueryClient) ReadTx(
	context.Context,
	*connect.Request[query.ReadTxRequest],
) (*connect.Response[query.ReadTxResponse], error) {
	return connect.NewResponse(&query.ReadTxResponse{}), nil
}

func (*recordingQueryClient) ReadGenesis(
	context.Context,
	*connect.Request[query.ReadGenesisRequest],
) (*connect.Response[query.ReadGenesisResponse], error) {
	return connect.NewResponse(&query.ReadGenesisResponse{}), nil
}

func (*recordingQueryClient) ReadEraSummary(
	context.Context,
	*connect.Request[query.ReadEraSummaryRequest],
) (*connect.Response[query.ReadEraSummaryResponse], error) {
	return connect.NewResponse(&query.ReadEraSummaryResponse{}), nil
}

func (*recordingQueryClient) ReadState(
	context.Context,
	*connect.Request[query.ReadStateRequest],
) (*connect.Response[query.ReadStateResponse], error) {
	return connect.NewResponse(&query.ReadStateResponse{}), nil
}

var _ sdk.QueryServiceClient = (*recordingQueryClient)(nil)
