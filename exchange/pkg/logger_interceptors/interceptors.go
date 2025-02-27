package loggerinterceptors

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func LoggingUnaryInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		md, _ := metadata.FromIncomingContext(ctx)
		clientAddr := "unknown"
		if vals := md.Get(":authority"); len(vals) > 0 {
			clientAddr = vals[0]
		}

		resp, err := handler(ctx, req)

		log.Info("gRPC Unary request",
			slog.String("method", info.FullMethod),
			slog.String("client", clientAddr),
			slog.Duration("duration", time.Since(start)),
			slog.Bool("success", err == nil),
		)

		return resp, err
	}
}

func LoggingStreamInterceptor(log *slog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		md, _ := metadata.FromIncomingContext(ss.Context())
		clientAddr := "unknown"
		if vals := md.Get(":authority"); len(vals) > 0 {
			clientAddr = vals[0]
		}

		err := handler(srv, ss)

		log.Info("gRPC Stream request",
			slog.String("method", info.FullMethod),
			slog.String("client", clientAddr),
			slog.Duration("duration", time.Since(start)),
			slog.Bool("success", err == nil),
		)

		return err
	}
}
