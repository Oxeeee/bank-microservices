package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	custerrors "github.com/Oxeeee/bank-microservices/billing/internal/models/errors"
	"github.com/Oxeeee/bank-microservices/billing/internal/models/requests"
	"github.com/Oxeeee/bank-microservices/billing/internal/models/responses"
	"github.com/Oxeeee/bank-microservices/billing/internal/service"
	jsonwrap "github.com/Oxeeee/bank-microservices/billing/pkg/request_json"
	reqvalidator "github.com/Oxeeee/bank-microservices/billing/pkg/request_validator"
)

type BillingHandler interface {
	Pay(w http.ResponseWriter, r *http.Request)
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

func (h *billingHandler) Pay(w http.ResponseWriter, r *http.Request) {
	const op = "handler.pay"
	log := h.log.With(slog.String("op", op))
	var req requests.BillPayment

	if err := jsonwrap.Unwrap(&req, r); err != nil {
		log.Info("bad request", "error", err)
		http.Error(w, "can not unmarshall json", http.StatusBadRequest)
		return
	}

	if err := reqvalidator.Validate(req); err != nil {
		http.Error(w, "has empty fields", http.StatusBadRequest)
		return
	}

	paymentID, err := h.service.Pay(&req)
	if err != nil {
		if err == custerrors.ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err == custerrors.ErrInsufficientBalance {
			http.Error(w, "insufficient balance", http.StatusPaymentRequired)
			return
		}

		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	var resp = responses.PaymentResponse{PaymentID: paymentID}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "error while encoding json", http.StatusOK)
		return
	}
}
