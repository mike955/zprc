package interceptor

import (
	"context"
	"runtime"
	"time"

	"github.com/mike955/zrpc/log"
	"google.golang.org/grpc"
)

func RecoveryInterceptor(logger *log.Entry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				logger.Errorf("recovery: %v: %+v\n%s\n", rerr, req, buf)
				// add err handle
			}
		}()
		return handler(ctx, req)
	}
}

func TimeoutInterceptor(logger *log.Entry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second*60)
		defer cancel()
		return handler(ctx, req)
	}
}
