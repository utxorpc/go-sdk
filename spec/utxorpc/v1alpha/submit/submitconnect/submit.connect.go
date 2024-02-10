// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: utxorpc/v1alpha/submit/submit.proto

package submitconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	submit "github.com/utxorpc/go-sdk/spec/utxorpc/v1alpha/submit"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// SubmitServiceName is the fully-qualified name of the SubmitService service.
	SubmitServiceName = "utxorpc.v1alpha.submit.SubmitService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// SubmitServiceSubmitTxProcedure is the fully-qualified name of the SubmitService's SubmitTx RPC.
	SubmitServiceSubmitTxProcedure = "/utxorpc.v1alpha.submit.SubmitService/SubmitTx"
	// SubmitServiceWaitForTxProcedure is the fully-qualified name of the SubmitService's WaitForTx RPC.
	SubmitServiceWaitForTxProcedure = "/utxorpc.v1alpha.submit.SubmitService/WaitForTx"
	// SubmitServiceReadMempoolProcedure is the fully-qualified name of the SubmitService's ReadMempool
	// RPC.
	SubmitServiceReadMempoolProcedure = "/utxorpc.v1alpha.submit.SubmitService/ReadMempool"
	// SubmitServiceWatchMempoolProcedure is the fully-qualified name of the SubmitService's
	// WatchMempool RPC.
	SubmitServiceWatchMempoolProcedure = "/utxorpc.v1alpha.submit.SubmitService/WatchMempool"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	submitServiceServiceDescriptor            = submit.File_utxorpc_v1alpha_submit_submit_proto.Services().ByName("SubmitService")
	submitServiceSubmitTxMethodDescriptor     = submitServiceServiceDescriptor.Methods().ByName("SubmitTx")
	submitServiceWaitForTxMethodDescriptor    = submitServiceServiceDescriptor.Methods().ByName("WaitForTx")
	submitServiceReadMempoolMethodDescriptor  = submitServiceServiceDescriptor.Methods().ByName("ReadMempool")
	submitServiceWatchMempoolMethodDescriptor = submitServiceServiceDescriptor.Methods().ByName("WatchMempool")
)

// SubmitServiceClient is a client for the utxorpc.v1alpha.submit.SubmitService service.
type SubmitServiceClient interface {
	SubmitTx(context.Context, *connect.Request[submit.SubmitTxRequest]) (*connect.Response[submit.SubmitTxResponse], error)
	WaitForTx(context.Context, *connect.Request[submit.WaitForTxRequest]) (*connect.ServerStreamForClient[submit.WaitForTxResponse], error)
	ReadMempool(context.Context, *connect.Request[submit.ReadMempoolRequest]) (*connect.Response[submit.ReadMempoolResponse], error)
	WatchMempool(context.Context, *connect.Request[submit.WatchMempoolRequest]) (*connect.ServerStreamForClient[submit.WatchMempoolResponse], error)
}

// NewSubmitServiceClient constructs a client for the utxorpc.v1alpha.submit.SubmitService service.
// By default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped
// responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewSubmitServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) SubmitServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &submitServiceClient{
		submitTx: connect.NewClient[submit.SubmitTxRequest, submit.SubmitTxResponse](
			httpClient,
			baseURL+SubmitServiceSubmitTxProcedure,
			connect.WithSchema(submitServiceSubmitTxMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		waitForTx: connect.NewClient[submit.WaitForTxRequest, submit.WaitForTxResponse](
			httpClient,
			baseURL+SubmitServiceWaitForTxProcedure,
			connect.WithSchema(submitServiceWaitForTxMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		readMempool: connect.NewClient[submit.ReadMempoolRequest, submit.ReadMempoolResponse](
			httpClient,
			baseURL+SubmitServiceReadMempoolProcedure,
			connect.WithSchema(submitServiceReadMempoolMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		watchMempool: connect.NewClient[submit.WatchMempoolRequest, submit.WatchMempoolResponse](
			httpClient,
			baseURL+SubmitServiceWatchMempoolProcedure,
			connect.WithSchema(submitServiceWatchMempoolMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// submitServiceClient implements SubmitServiceClient.
type submitServiceClient struct {
	submitTx     *connect.Client[submit.SubmitTxRequest, submit.SubmitTxResponse]
	waitForTx    *connect.Client[submit.WaitForTxRequest, submit.WaitForTxResponse]
	readMempool  *connect.Client[submit.ReadMempoolRequest, submit.ReadMempoolResponse]
	watchMempool *connect.Client[submit.WatchMempoolRequest, submit.WatchMempoolResponse]
}

// SubmitTx calls utxorpc.v1alpha.submit.SubmitService.SubmitTx.
func (c *submitServiceClient) SubmitTx(ctx context.Context, req *connect.Request[submit.SubmitTxRequest]) (*connect.Response[submit.SubmitTxResponse], error) {
	return c.submitTx.CallUnary(ctx, req)
}

// WaitForTx calls utxorpc.v1alpha.submit.SubmitService.WaitForTx.
func (c *submitServiceClient) WaitForTx(ctx context.Context, req *connect.Request[submit.WaitForTxRequest]) (*connect.ServerStreamForClient[submit.WaitForTxResponse], error) {
	return c.waitForTx.CallServerStream(ctx, req)
}

// ReadMempool calls utxorpc.v1alpha.submit.SubmitService.ReadMempool.
func (c *submitServiceClient) ReadMempool(ctx context.Context, req *connect.Request[submit.ReadMempoolRequest]) (*connect.Response[submit.ReadMempoolResponse], error) {
	return c.readMempool.CallUnary(ctx, req)
}

// WatchMempool calls utxorpc.v1alpha.submit.SubmitService.WatchMempool.
func (c *submitServiceClient) WatchMempool(ctx context.Context, req *connect.Request[submit.WatchMempoolRequest]) (*connect.ServerStreamForClient[submit.WatchMempoolResponse], error) {
	return c.watchMempool.CallServerStream(ctx, req)
}

// SubmitServiceHandler is an implementation of the utxorpc.v1alpha.submit.SubmitService service.
type SubmitServiceHandler interface {
	SubmitTx(context.Context, *connect.Request[submit.SubmitTxRequest]) (*connect.Response[submit.SubmitTxResponse], error)
	WaitForTx(context.Context, *connect.Request[submit.WaitForTxRequest], *connect.ServerStream[submit.WaitForTxResponse]) error
	ReadMempool(context.Context, *connect.Request[submit.ReadMempoolRequest]) (*connect.Response[submit.ReadMempoolResponse], error)
	WatchMempool(context.Context, *connect.Request[submit.WatchMempoolRequest], *connect.ServerStream[submit.WatchMempoolResponse]) error
}

// NewSubmitServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewSubmitServiceHandler(svc SubmitServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	submitServiceSubmitTxHandler := connect.NewUnaryHandler(
		SubmitServiceSubmitTxProcedure,
		svc.SubmitTx,
		connect.WithSchema(submitServiceSubmitTxMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	submitServiceWaitForTxHandler := connect.NewServerStreamHandler(
		SubmitServiceWaitForTxProcedure,
		svc.WaitForTx,
		connect.WithSchema(submitServiceWaitForTxMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	submitServiceReadMempoolHandler := connect.NewUnaryHandler(
		SubmitServiceReadMempoolProcedure,
		svc.ReadMempool,
		connect.WithSchema(submitServiceReadMempoolMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	submitServiceWatchMempoolHandler := connect.NewServerStreamHandler(
		SubmitServiceWatchMempoolProcedure,
		svc.WatchMempool,
		connect.WithSchema(submitServiceWatchMempoolMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/utxorpc.v1alpha.submit.SubmitService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case SubmitServiceSubmitTxProcedure:
			submitServiceSubmitTxHandler.ServeHTTP(w, r)
		case SubmitServiceWaitForTxProcedure:
			submitServiceWaitForTxHandler.ServeHTTP(w, r)
		case SubmitServiceReadMempoolProcedure:
			submitServiceReadMempoolHandler.ServeHTTP(w, r)
		case SubmitServiceWatchMempoolProcedure:
			submitServiceWatchMempoolHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedSubmitServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedSubmitServiceHandler struct{}

func (UnimplementedSubmitServiceHandler) SubmitTx(context.Context, *connect.Request[submit.SubmitTxRequest]) (*connect.Response[submit.SubmitTxResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("utxorpc.v1alpha.submit.SubmitService.SubmitTx is not implemented"))
}

func (UnimplementedSubmitServiceHandler) WaitForTx(context.Context, *connect.Request[submit.WaitForTxRequest], *connect.ServerStream[submit.WaitForTxResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("utxorpc.v1alpha.submit.SubmitService.WaitForTx is not implemented"))
}

func (UnimplementedSubmitServiceHandler) ReadMempool(context.Context, *connect.Request[submit.ReadMempoolRequest]) (*connect.Response[submit.ReadMempoolResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("utxorpc.v1alpha.submit.SubmitService.ReadMempool is not implemented"))
}

func (UnimplementedSubmitServiceHandler) WatchMempool(context.Context, *connect.Request[submit.WatchMempoolRequest], *connect.ServerStream[submit.WatchMempoolResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("utxorpc.v1alpha.submit.SubmitService.WatchMempool is not implemented"))
}
