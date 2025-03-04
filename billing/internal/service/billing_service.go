package service

import (
	"log/slog"

	"github.com/Oxeeee/bank-microservices/billing/internal/config"
	"github.com/Oxeeee/bank-microservices/billing/internal/models/domain"
	"github.com/Oxeeee/bank-microservices/billing/internal/models/requests"
	"github.com/Oxeeee/bank-microservices/billing/internal/repo"
	"github.com/google/uuid"
)

type BillingService interface {
	// GetUserByID — takes uuid and return model of user and error
	GetUserByID(uuid uuid.UUID) (*domain.User, error)

	Pay(req *requests.BillPayment) (uuid.UUID, error)
}

type billingService struct {
	log   *slog.Logger
	cfg   *config.Config
	repo  repo.BillingRepository
	cache repo.BillingCache
	kafka repo.BillingKafkaRepo
}

// NewBillingService create new example of billingService structure
func NewBillingService(log *slog.Logger, cfg *config.Config, repo repo.BillingRepository, cache repo.BillingCache, kafka repo.BillingKafkaRepo) BillingService {
	return &billingService{
		log:   log,
		cfg:   cfg,
		repo:  repo,
		cache: cache,
		kafka: kafka,
	}
}

// GetUserByID — takes uuid and return model of user and error
func (s *billingService) GetUserByID(uuid uuid.UUID) (*domain.User, error) {
	user, err := s.repo.GetUserByID(uuid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *billingService) Pay(req *requests.BillPayment) (uuid.UUID, error) {
	return s.repo.ProcessPayment(req)
}
