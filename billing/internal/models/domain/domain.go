package domain

import (
	"time"

	"encoding/json"

	"github.com/google/uuid"
)

type BillPayment struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Provider  string          `json:"provider"`
	Amount    float64         `json:"amount"`
	Status    string          `json:"status"`
	Details   json.RawMessage `json:"details,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
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
