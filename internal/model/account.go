package model

import (
	"github.com/google/uuid"
)

type WalletBalance struct {
	ID     uuid.UUID `json:"-"`
	Amount float64   `json:"amount"`
	UserID uuid.UUID `json:"user_id"`
}
