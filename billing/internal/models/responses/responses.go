package responses

import "github.com/google/uuid"

type PaymentResponse struct {
	PaymentID uuid.UUID `json:"payment_id"`
}
