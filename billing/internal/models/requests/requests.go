package requests

import (
	"encoding/json"

	"github.com/google/uuid"
)

type BillPayment struct {
	UserID   uuid.UUID       `json:"user_id"`
	Provider string          `json:"provider"`
	Amount   float64         `json:"amount"`
	Currency string          `json:"currency"`
	Details  json.RawMessage `json:"details,omitempty"`
}
