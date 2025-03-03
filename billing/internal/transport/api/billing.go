package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/Oxeeee/bank-microservices/billing/internal/models/requests"
	"github.com/Oxeeee/bank-microservices/billing/internal/service"
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
	// TODO: вместо прямого запроса из нашего микросервиса — поход в auth за юзером (зарефачь все слои и удали нах все че не нужно)
	var req requests.BillPayment

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error while reading request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error while parse JSON", http.StatusBadRequest)
		return
	}

	if err := reqvalidator.Validate(req); err != nil {
		http.Error(w, "has empty fields", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByID(req.UserID)
	if err != nil {
		http.Error(w, "error while get user by id", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "error while encoding json", http.StatusInternalServerError)
		return
	}
}
