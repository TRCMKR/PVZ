package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"gitlab.ozon.dev/alexplay1224/homework/pkg/monitoring"
)

// MetricsInterceptor is an interceptor that updates metrics
func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		monitoring.SetRequestCounter()
		monitoring.SetGrpcRequestByMethodCount(info.FullMethod)

		resp, err := handler(ctx, req)

		duration := time.Since(start).Seconds()
		monitoring.RequestDuration(duration)
		monitoring.SetResponseTimeSummary(duration)

		if err != nil {
			monitoring.SetErrorCounter()
			monitoring.SetGrpcRequestCountWithStatus("error")
		} else {
			monitoring.SetGrpcRequestCountWithStatus("ok")
		}

		return resp, err
	}
}
