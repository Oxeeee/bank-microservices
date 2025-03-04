package domain

import (
	"time"

	"encoding/json"

	"github.com/google/uuid"
)

type BillPayment struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	UserID    uuid.UUID       `json:"user_id" db:"user_id"`
	Provider  string          `json:"provider" db:"provider"`
	Amount    float64         `json:"amount" db:"amount"`
	Currency  string          `json:"currency" db:"currency"` // добавлено поле валюты
	Status    string          `json:"status" db:"status"`
	Details   json.RawMessage `json:"details,omitempty" db:"details"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
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
