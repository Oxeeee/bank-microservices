package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Oxeeee/bank-microservices/billing/internal/service"
)

type BillingHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
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

type Response struct {
	Message string
}

func (h *billingHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	resp := Response{Message: "Hello world!"}
	json.NewEncoder(w).Encode(resp)
}
