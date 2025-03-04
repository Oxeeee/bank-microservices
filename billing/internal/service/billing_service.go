package service

import (
	"log/slog"

	"github.com/Oxeeee/bank-microservices/billing/internal/config"

	"github.com/Oxeeee/bank-microservices/billing/internal/models/requests"
	"github.com/Oxeeee/bank-microservices/billing/internal/repo"
	"github.com/google/uuid"
)

type BillingService interface {
	// GetUserByID — takes uuid and return model of user and error

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

func (s *billingService) Pay(req *requests.BillPayment) (uuid.UUID, error) {
	paymentID, err := s.repo.ProcessPayment(req)
	if err != nil {
		return uuid.Nil, err
	}

	bill, err := s.repo.GetPaymentByID(paymentID)
	if err != nil {
		s.log.Error("error while get payment by id", "error", err)
		return paymentID, nil
	}

	err = s.kafka.SendPaymentStatus(bill)
	if err != nil {
		s.log.Error("error while send message to kafka", "error", err)
		return paymentID, nil
	}
	return paymentID, nil
}
