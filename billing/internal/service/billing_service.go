package service

import (
	"log/slog"

	"github.com/Oxeeee/bank-microservices/billing/internal/config"
	"github.com/Oxeeee/bank-microservices/billing/internal/repo"
)

type BillingService interface {
}

type billingService struct {
	log   *slog.Logger
	cfg   *config.Config
	repo  repo.BillingRepository
	cache repo.BillingCache
}

func NewBillingService(log *slog.Logger, cfg *config.Config, repo repo.BillingRepository, cache repo.BillingCache) BillingService {
	return &billingService{
		log:   log,
		cfg:   cfg,
		repo:  repo,
		cache: cache,
	}
}

