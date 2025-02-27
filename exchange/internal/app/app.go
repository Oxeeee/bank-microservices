package app

import (
	"log/slog"

	grpcapp "github.com/Oxeeee/bank-microservices/exchange/internal/app/grpc"
	"github.com/Oxeeee/bank-microservices/exchange/internal/service"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, service service.ExchangeService) *App {
	grpcapp := grpcapp.New(log, grpcPort, service)

	return &App{
		GRPCSrv: grpcapp,
	}
}
