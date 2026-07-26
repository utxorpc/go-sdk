package sdk

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1beta/query"
	"github.com/utxorpc/go-codegen/utxorpc/v1beta/sync"
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

func TestReadStateWrappersAddHeadersAndCallQueryClient(t *testing.T) {
	tests := []struct {
		name string
		call func(*UtxorpcClient, *connect.Request[query.ReadStateRequest]) error
	}{
		{
			name: "background context",
			call: func(
				client *UtxorpcClient,
				req *connect.Request[query.ReadStateRequest],
			) error {
				_, err := client.ReadState(req)
				return err
			},
		},
		{
			name: "provided context",
			call: func(
				client *UtxorpcClient,
				req *connect.Request[query.ReadStateRequest],
			) error {
				_, err := client.ReadStateWithContext(context.Background(), req)
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fakeQuery := &recordingQueryClient{}
			client := NewClient(
				WithHeaders(map[string]string{"dmtr-api-key": "secret"}),
			)
			client.Query = fakeQuery

			if err := test.call(
				client,
				connect.NewRequest(&query.ReadStateRequest{}),
			); err != nil {
				t.Fatalf("ReadState wrapper returned error: %v", err)
			}
			if fakeQuery.readStateReq == nil {
				t.Fatal("Query.ReadState was not called")
			}
			if got := fakeQuery.readStateReq.Header().
				Get("dmtr-api-key"); got != "secret" {
				t.Fatalf("dmtr-api-key header = %q, want %q", got, "secret")
			}
		})
	}
}

func TestDumpHistoryWrappersAddHeadersAndCallSyncClient(t *testing.T) {
	tests := []struct {
		name string
		call func(*UtxorpcClient, *connect.Request[sync.DumpHistoryRequest]) error
	}{
		{
			name: "background context",
			call: func(
				client *UtxorpcClient,
				req *connect.Request[sync.DumpHistoryRequest],
			) error {
				_, err := client.DumpHistory(req)
				return err
			},
		},
		{
			name: "provided context",
			call: func(
				client *UtxorpcClient,
				req *connect.Request[sync.DumpHistoryRequest],
			) error {
				_, err := client.DumpHistoryWithContext(context.Background(), req)
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fakeSync := &recordingSyncClient{}
			client := NewClient(
				WithHeaders(map[string]string{"dmtr-api-key": "secret"}),
			)
			client.Sync = fakeSync

			if err := test.call(
				client,
				connect.NewRequest(&sync.DumpHistoryRequest{}),
			); err != nil {
				t.Fatalf("DumpHistory wrapper returned error: %v", err)
			}
			if fakeSync.dumpHistoryReq == nil {
				t.Fatal("Sync.DumpHistory was not called")
			}
			if got := fakeSync.dumpHistoryReq.Header().
				Get("dmtr-api-key"); got != "secret" {
				t.Fatalf("dmtr-api-key header = %q, want %q", got, "secret")
			}
		})
	}
}

type recordingQueryClient struct {
	readParamsReq *connect.Request[query.ReadParamsRequest]
	readStateReq  *connect.Request[query.ReadStateRequest]
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

func (r *recordingQueryClient) ReadState(
	_ context.Context,
	req *connect.Request[query.ReadStateRequest],
) (*connect.Response[query.ReadStateResponse], error) {
	r.readStateReq = req
	return connect.NewResponse(&query.ReadStateResponse{}), nil
}

var _ QueryServiceClient = (*recordingQueryClient)(nil)

type recordingSyncClient struct {
	dumpHistoryReq *connect.Request[sync.DumpHistoryRequest]
}

func (*recordingSyncClient) FetchBlock(
	context.Context,
	*connect.Request[sync.FetchBlockRequest],
) (*connect.Response[sync.FetchBlockResponse], error) {
	return connect.NewResponse(&sync.FetchBlockResponse{}), nil
}

func (r *recordingSyncClient) DumpHistory(
	_ context.Context,
	req *connect.Request[sync.DumpHistoryRequest],
) (*connect.Response[sync.DumpHistoryResponse], error) {
	r.dumpHistoryReq = req
	return connect.NewResponse(&sync.DumpHistoryResponse{}), nil
}

func (*recordingSyncClient) FollowTip(
	context.Context,
	*connect.Request[sync.FollowTipRequest],
) (*connect.ServerStreamForClient[sync.FollowTipResponse], error) {
	return nil, nil
}

func (*recordingSyncClient) ReadTip(
	context.Context,
	*connect.Request[sync.ReadTipRequest],
) (*connect.Response[sync.ReadTipResponse], error) {
	return connect.NewResponse(&sync.ReadTipResponse{}), nil
}

var _ SyncServiceClient = (*recordingSyncClient)(nil)
