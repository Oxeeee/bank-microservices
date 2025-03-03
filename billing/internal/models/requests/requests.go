package requests

import (
	"github.com/google/uuid"
)

type BillPayment struct {
	UserID        uuid.UUID `json:"user_id"`
	Provider      string    `json:"provider"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
}
