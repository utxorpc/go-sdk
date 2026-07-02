package sdk

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1beta/query"
)

func TestHeaderManagementAndRequestInjection(t *testing.T) {
	client := NewClient()

	client.SetHeader("dmtr-api-key", "secret")
	client.SetHeader("x-extra", "value")
	client.RemoveHeader("x-extra")

	req := connect.NewRequest(&query.ReadParamsRequest{})
	client.AddHeadersToRequest(req)

	if got := req.Header().Get("dmtr-api-key"); got != "secret" {
		t.Fatalf("dmtr-api-key header = %q, want %q", got, "secret")
	}
	if got := req.Header().Get("x-extra"); got != "" {
		t.Fatalf("removed header x-extra = %q, want empty", got)
	}

	client.SetHeaders(map[string]string{"authorization": "Bearer token"})
	headers := client.Headers()
	if got := headers["authorization"]; got != "Bearer token" {
		t.Fatalf("authorization header = %q, want %q", got, "Bearer token")
	}
}

func TestSetURLRebuildsServiceClients(t *testing.T) {
	client := NewClient(WithBaseUrl("http://old.example.test"))

	oldQuery := client.Query
	oldSubmit := client.Submit
	oldSync := client.Sync
	oldWatch := client.Watch

	client.SetURL("http://new.example.test")

	if got := client.URL(); got != "http://new.example.test" {
		t.Fatalf("URL() = %q, want %q", got, "http://new.example.test")
	}
	if client.Query == oldQuery {
		t.Fatal("Query service client was not rebuilt")
	}
	if client.Submit == oldSubmit {
		t.Fatal("Submit service client was not rebuilt")
	}
	if client.Sync == oldSync {
		t.Fatal("Sync service client was not rebuilt")
	}
	if client.Watch == oldWatch {
		t.Fatal("Watch service client was not rebuilt")
	}
}

func TestReadParamsWithContextAddsHeaders(t *testing.T) {
	fakeQuery := &recordingQueryClient{}
	client := NewClient(
		WithBaseUrl("http://example.test"),
		WithHeaders(map[string]string{"dmtr-api-key": "secret"}),
	)
	client.Query = fakeQuery

	resp, err := client.ReadParamsWithContext(
		context.Background(),
		connect.NewRequest(&query.ReadParamsRequest{}),
	)
	if err != nil {
		t.Fatalf("ReadParamsWithContext returned error: %v", err)
	}
	if resp == nil {
		t.Fatal("ReadParamsWithContext returned nil response")
	}
	if fakeQuery.readParamsReq == nil {
		t.Fatal("ReadParams was not called")
	}
	if got := fakeQuery.readParamsReq.Header().Get("dmtr-api-key"); got != "secret" {
		t.Fatalf("dmtr-api-key header = %q, want %q", got, "secret")
	}
}

type recordingQueryClient struct {
	readParamsReq *connect.Request[query.ReadParamsRequest]
}

func (r *recordingQueryClient) ReadParams(
	_ context.Context,
	req *connect.Request[query.ReadParamsRequest],
) (*connect.Response[query.ReadParamsResponse], error) {
	r.readParamsReq = req
	return connect.NewResponse(&query.ReadParamsResponse{}), nil
}

func (*recordingQueryClient) ReadUtxos(
	context.Context,
	*connect.Request[query.ReadUtxosRequest],
) (*connect.Response[query.ReadUtxosResponse], error) {
	return connect.NewResponse(&query.ReadUtxosResponse{}), nil
}

func (*recordingQueryClient) SearchUtxos(
	context.Context,
	*connect.Request[query.SearchUtxosRequest],
) (*connect.Response[query.SearchUtxosResponse], error) {
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

var _ QueryServiceClient = (*recordingQueryClient)(nil)
