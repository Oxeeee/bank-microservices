package main

import (
	"log/slog"
	"os"

	"github.com/Oxeeee/bank-microservices/billing/internal/config"
	"github.com/Oxeeee/bank-microservices/billing/internal/db"
	"github.com/Oxeeee/bank-microservices/billing/internal/repo"
	"github.com/Oxeeee/bank-microservices/billing/internal/service"
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

	database := db.InitDB(cfg)
	redis := db.InitRedis(cfg.Redis.Address, cfg.Redis.Password, cfg.Redis.DB)

	cacheRepo := repo.NewBillingCache(redis)
	dbRepo := repo.NewBillingRepository(database)

	service := service.NewBillingService(log, cfg, dbRepo, cacheRepo)
	
	
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
