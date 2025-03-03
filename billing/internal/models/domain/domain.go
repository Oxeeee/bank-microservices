package domain

import (
	"time"

	"encoding/json"

	"github.com/google/uuid"
)

type BillPayment struct {
	ID            uuid.UUID       `json:"id"`
	UserID        uuid.UUID       `json:"user_id"`
	Provider      string          `json:"provider"`
	Amount        float64         `json:"amount"`
	Currency      string          `json:"currency"`
	PaymentMethod string          `json:"payment_method"`
	Status        string          `json:"status"`
	Details       json.RawMessage `json:"details,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type Provider struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Type   string    `json:"type"`
	Active bool      `json:"active"`
}

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash,omitempty" db:"password_hash"`
	Balance      float64   `json:"balance" db:"balance"`
	Currency     string    `json:"currency" db:"currency"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
