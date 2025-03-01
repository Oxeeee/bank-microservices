package api

import (
	"log/slog"
	"net/http"

	"github.com/Oxeeee/bank-microservices/billing/internal/service"
)

type BillingHandler interface {
}

type billingHandler struct {
	log     *slog.Logger
	service service.BillingService
}

func NewBillingHandler(log *slog.Logger, service service.BillingService) BillingHandler {
	return &billingHandler{
		log:     log,
		service: service,
	}
}

func (h *billingHandler) Register(w http.ResponseWriter, r *http.Request) {

}
