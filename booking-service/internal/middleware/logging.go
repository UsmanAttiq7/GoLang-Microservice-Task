package middleware

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func LoggingInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Log incoming request
		logger.WithFields(logrus.Fields{
			"method":  info.FullMethod,
			"request": req,
		}).Info("gRPC Request")

		// Handle the request
		resp, err := handler(ctx, req)

		// Log response or error
		if err != nil {
			logger.WithFields(logrus.Fields{
				"method": info.FullMethod,
				"error":  err.Error(),
			}).Error("gRPC Response Error")
		} else {
			logger.WithFields(logrus.Fields{
				"method":   info.FullMethod,
				"response": resp,
			}).Info("gRPC Response")
		}

		return resp, err
	}
}
