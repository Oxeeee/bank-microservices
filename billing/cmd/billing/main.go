package main

import (
	"log/slog"
	"os"

	"github.com/Oxeeee/bank-microservices/billing/internal/app"
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
	log.Info("connected to database")

	redis := db.InitRedis(cfg.Redis.Address, cfg.Redis.Password, cfg.Redis.DB)
	log.Info("connected to redis")

	producer := db.InitKafka(cfg.Kafka.Brokers)
	log.Info("connected to kafka")

	dbRepo := repo.NewBillingRepository(database, log)
	log.Info("initialized dbRepo")

	cacheRepo := repo.NewBillingCache(redis)
	log.Info("initialized redisRepo")

	kafkaRepo := repo.NewBillingKafkaRepo(producer, cfg.Kafka.Topic, log)

	service := service.NewBillingService(log, cfg, dbRepo, cacheRepo, kafkaRepo)
	log.Info("initialized service")

	application := app.New(log, cfg.RESTPort, service)
	application.RESTSrv.MustRun()
	log.Info("started")
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
