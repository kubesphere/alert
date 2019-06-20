// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: custom.proto

/*
Package pb is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package pb

import (
	"context"
	"io"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray

var (
	filter_AlertManagerCustom_DescribeAlertsWithResource_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_AlertManagerCustom_DescribeAlertsWithResource_0(ctx context.Context, marshaler runtime.Marshaler, client AlertManagerCustomClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq DescribeAlertsWithResourceRequest
	var metadata runtime.ServerMetadata

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_AlertManagerCustom_DescribeAlertsWithResource_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.DescribeAlertsWithResource(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

var (
	filter_AlertManagerCustom_DescribeAlertDetails_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_AlertManagerCustom_DescribeAlertDetails_0(ctx context.Context, marshaler runtime.Marshaler, client AlertManagerCustomClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq DescribeAlertDetailsRequest
	var metadata runtime.ServerMetadata

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_AlertManagerCustom_DescribeAlertDetails_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.DescribeAlertDetails(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

var (
	filter_AlertManagerCustom_DescribeAlertStatus_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_AlertManagerCustom_DescribeAlertStatus_0(ctx context.Context, marshaler runtime.Marshaler, client AlertManagerCustomClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq DescribeAlertStatusRequest
	var metadata runtime.ServerMetadata

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_AlertManagerCustom_DescribeAlertStatus_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.DescribeAlertStatus(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

var (
	filter_AlertManagerCustom_DescribeHistoryDetail_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_AlertManagerCustom_DescribeHistoryDetail_0(ctx context.Context, marshaler runtime.Marshaler, client AlertManagerCustomClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq DescribeHistoryDetailRequest
	var metadata runtime.ServerMetadata

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_AlertManagerCustom_DescribeHistoryDetail_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.DescribeHistoryDetail(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

// RegisterAlertManagerCustomHandlerFromEndpoint is same as RegisterAlertManagerCustomHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterAlertManagerCustomHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterAlertManagerCustomHandler(ctx, mux, conn)
}

// RegisterAlertManagerCustomHandler registers the http handlers for service AlertManagerCustom to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterAlertManagerCustomHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterAlertManagerCustomHandlerClient(ctx, mux, NewAlertManagerCustomClient(conn))
}

// RegisterAlertManagerCustomHandlerClient registers the http handlers for service AlertManagerCustom
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "AlertManagerCustomClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "AlertManagerCustomClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "AlertManagerCustomClient" to call the correct interceptors.
func RegisterAlertManagerCustomHandlerClient(ctx context.Context, mux *runtime.ServeMux, client AlertManagerCustomClient) error {

	mux.Handle("GET", pattern_AlertManagerCustom_DescribeAlertsWithResource_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_AlertManagerCustom_DescribeAlertsWithResource_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AlertManagerCustom_DescribeAlertsWithResource_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_AlertManagerCustom_DescribeAlertDetails_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_AlertManagerCustom_DescribeAlertDetails_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AlertManagerCustom_DescribeAlertDetails_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_AlertManagerCustom_DescribeAlertStatus_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_AlertManagerCustom_DescribeAlertStatus_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AlertManagerCustom_DescribeAlertStatus_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_AlertManagerCustom_DescribeHistoryDetail_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_AlertManagerCustom_DescribeHistoryDetail_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AlertManagerCustom_DescribeHistoryDetail_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_AlertManagerCustom_DescribeAlertsWithResource_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"v1", "alert_with_resource"}, ""))

	pattern_AlertManagerCustom_DescribeAlertDetails_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"v1", "alert_detail"}, ""))

	pattern_AlertManagerCustom_DescribeAlertStatus_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"v1", "alert_status"}, ""))

	pattern_AlertManagerCustom_DescribeHistoryDetail_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"v1", "historydetail"}, ""))
)

var (
	forward_AlertManagerCustom_DescribeAlertsWithResource_0 = runtime.ForwardResponseMessage

	forward_AlertManagerCustom_DescribeAlertDetails_0 = runtime.ForwardResponseMessage

	forward_AlertManagerCustom_DescribeAlertStatus_0 = runtime.ForwardResponseMessage

	forward_AlertManagerCustom_DescribeHistoryDetail_0 = runtime.ForwardResponseMessage
)
