package grpc

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ClientUnaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "traceId", ctx.Value("traceId").(string))
	ctx = metadata.AppendToOutgoingContext(ctx, "x_real_ip", ctx.Value("x_real_ip").(string))
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	end := time.Now()
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Infof("info: rpc call,method: %s start time: %s, end time: %s, err: %v", method, start.Format("Basic"), end.Format(time.RFC3339), err)
	return err
}

func ClientStreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "traceId", ctx.Value("traceId").(string))
	ctx = metadata.AppendToOutgoingContext(ctx, "x_real_ip", ctx.Value("x_real_ip").(string))
	logger := ctx.Value("logger").(*logrus.Entry)
	start := time.Now()
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return newWrappedStream(ctx, s, start, logger), nil
}

type wrappedStream struct {
	ctx context.Context
	grpc.ClientStream
	start  time.Time
	logger *logrus.Entry
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	w.logger.Infof("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	w.logger.Infof("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

func newWrappedStream(ctx context.Context, s grpc.ClientStream, start time.Time, logger *logrus.Entry) grpc.ClientStream {
	return &wrappedStream{ctx, s, start, logger}
}
