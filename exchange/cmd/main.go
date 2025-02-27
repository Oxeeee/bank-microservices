package main

import (
	"log/slog"
	"os"

	"github.com/Oxeeee/bank-microservices/exchange/internal/app"
	"github.com/Oxeeee/bank-microservices/exchange/internal/config"
	"github.com/Oxeeee/bank-microservices/exchange/internal/service"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting application")

	service := service.NewExchangeService(log, cfg)

	app := app.New(log, cfg.GRPC.Port, service)
	app.GRPCSrv.MustRun()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
