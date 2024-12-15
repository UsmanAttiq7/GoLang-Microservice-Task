package middleware

import (
	"context"
	"github.com/golang_falcon_task/booking-service/internal/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// MetricsInterceptor captures Prometheus metrics for gRPC calls.
func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Handle request
		resp, err := handler(ctx, req)

		// Update metrics
		st, _ := status.FromError(err)
		metrics.RequestCount.WithLabelValues(info.FullMethod, st.Code().String()).Inc()

		return resp, err
	}
}
