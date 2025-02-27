package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/Oxeeee/bank-microservices/exchange/internal/service"
	grpctransport "github.com/Oxeeee/bank-microservices/exchange/internal/transport/grpc"
	loggerinterceptors "github.com/Oxeeee/bank-microservices/exchange/pkg/logger_interceptors"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, exchangeService service.ExchangeService) *App {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggerinterceptors.LoggingUnaryInterceptor(log)),
		grpc.StreamInterceptor(loggerinterceptors.LoggingStreamInterceptor(log)),
	)

	grpctransport.Register(gRPCServer, exchangeService, log)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	const op = "grpcapp.run"
	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		panic(err)
	}

	log.Info("gRPC server is running")

	if err := a.gRPCServer.Serve(l); err != nil {
		panic(fmt.Sprintf("%s: %v", op, err))
	}
}
