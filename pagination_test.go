package sdk

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1beta/query"
	"github.com/utxorpc/go-codegen/utxorpc/v1beta/sync"
	"google.golang.org/protobuf/proto"
)

func TestSearchUtxosPagesFollowsTokensAndYieldsErrors(t *testing.T) {
	stopErr := errors.New("page failed")
	fakeQuery := &paginatedQueryClient{
		recordingQueryClient: &recordingQueryClient{},
		responses: []*query.SearchUtxosResponse{
			{NextToken: proto.String("second")},
			{NextToken: proto.String("third")},
		},
		err: stopErr,
	}
	client := NewClient(
		WithBaseUrl("http://example.test"),
		WithHeaders(map[string]string{"client-header": "client"}),
	)
	client.Query = fakeQuery

	req := connect.NewRequest(&query.SearchUtxosRequest{
		StartToken: proto.String("first"),
	})
	req.Header().Set("request-header", "request")

	var pageCount int
	var gotErr error
	for resp, err := range client.SearchUtxosPages(req) {
		if err != nil {
			gotErr = err
			continue
		}
		if resp == nil {
			t.Fatal("pagination yielded a nil response without an error")
		}
		pageCount++
	}

	if pageCount != 2 {
		t.Fatalf("page count = %d, want 2", pageCount)
	}
	if !errors.Is(gotErr, stopErr) {
		t.Fatalf("pagination error = %v, want %v", gotErr, stopErr)
	}
	wantTokens := []string{"first", "second", "third"}
	if len(fakeQuery.tokens) != len(wantTokens) {
		t.Fatalf("tokens = %q, want %q", fakeQuery.tokens, wantTokens)
	}
	for i, want := range wantTokens {
		if fakeQuery.tokens[i] != want {
			t.Fatalf("tokens[%d] = %q, want %q", i, fakeQuery.tokens[i], want)
		}
	}
	if fakeQuery.requestHeaders != len(wantTokens) {
		t.Fatalf(
			"requests with request-header = %d, want %d",
			fakeQuery.requestHeaders,
			len(wantTokens),
		)
	}
	if fakeQuery.clientHeaders != len(wantTokens) {
		t.Fatalf(
			"requests with client-header = %d, want %d",
			fakeQuery.clientHeaders,
			len(wantTokens),
		)
	}
	if req.Msg.GetStartToken() != "first" {
		t.Fatalf(
			"original start_token = %q, want %q",
			req.Msg.GetStartToken(),
			"first",
		)
	}
}

func TestSearchUtxosPagesStopsWhenCallerStops(t *testing.T) {
	fakeQuery := &paginatedQueryClient{
		recordingQueryClient: &recordingQueryClient{},
		responses: []*query.SearchUtxosResponse{
			{NextToken: proto.String("second")},
			{},
		},
	}
	client := NewClient(WithBaseUrl("http://example.test"))
	client.Query = fakeQuery

	for range client.SearchUtxosPages(
		connect.NewRequest(&query.SearchUtxosRequest{}),
	) {
		break
	}

	if len(fakeQuery.tokens) != 1 {
		t.Fatalf("request count = %d, want 1", len(fakeQuery.tokens))
	}
}

func TestDumpHistoryPagesFollowsBlockRefTokens(t *testing.T) {
	fakeSync := &paginatedSyncClient{
		recordingSyncClient: &recordingSyncClient{},
		responses: []*sync.DumpHistoryResponse{
			{NextToken: &sync.BlockRef{Slot: 20}},
			{},
		},
	}
	client := NewClient(WithBaseUrl("http://example.test"))
	client.Sync = fakeSync

	req := connect.NewRequest(&sync.DumpHistoryRequest{
		StartToken: &sync.BlockRef{Slot: 10},
	})

	var pageCount int
	for _, err := range client.DumpHistoryPagesWithContext(
		context.Background(),
		req,
	) {
		if err != nil {
			t.Fatalf("DumpHistoryPagesWithContext returned error: %v", err)
		}
		pageCount++
	}

	if pageCount != 2 {
		t.Fatalf("page count = %d, want 2", pageCount)
	}
	wantSlots := []uint64{10, 20}
	if len(fakeSync.slots) != len(wantSlots) {
		t.Fatalf("slots = %v, want %v", fakeSync.slots, wantSlots)
	}
	for i, want := range wantSlots {
		if fakeSync.slots[i] != want {
			t.Fatalf("slots[%d] = %d, want %d", i, fakeSync.slots[i], want)
		}
	}
	if req.Msg.GetStartToken().GetSlot() != 10 {
		t.Fatalf(
			"original start token slot = %d, want 10",
			req.Msg.GetStartToken().GetSlot(),
		)
	}
}

type paginatedQueryClient struct {
	*recordingQueryClient
	responses []*query.SearchUtxosResponse
	tokens    []string
	err       error

	requestHeaders int
	clientHeaders  int
}

func (p *paginatedQueryClient) SearchUtxos(
	_ context.Context,
	req *connect.Request[query.SearchUtxosRequest],
) (*connect.Response[query.SearchUtxosResponse], error) {
	p.tokens = append(p.tokens, req.Msg.GetStartToken())
	if req.Header().Get("request-header") == "request" {
		p.requestHeaders++
	}
	if req.Header().Get("client-header") == "client" {
		p.clientHeaders++
	}
	if len(p.responses) == 0 {
		return nil, p.err
	}
	resp := p.responses[0]
	p.responses = p.responses[1:]
	return connect.NewResponse(resp), nil
}

type paginatedSyncClient struct {
	*recordingSyncClient
	responses []*sync.DumpHistoryResponse
	slots     []uint64
}

func (p *paginatedSyncClient) DumpHistory(
	_ context.Context,
	req *connect.Request[sync.DumpHistoryRequest],
) (*connect.Response[sync.DumpHistoryResponse], error) {
	p.slots = append(p.slots, req.Msg.GetStartToken().GetSlot())
	resp := p.responses[0]
	p.responses = p.responses[1:]
	return connect.NewResponse(resp), nil
}

var (
	_ QueryServiceClient = (*paginatedQueryClient)(nil)
	_ SyncServiceClient  = (*paginatedSyncClient)(nil)
)
