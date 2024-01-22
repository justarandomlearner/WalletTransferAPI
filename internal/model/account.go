package model

import (
	"github.com/google/uuid"
)

type AccountBalance struct {
	ID        uuid.UUID `json:"-"`
	Amount    float64   `json:"amount"`
	AccountID uuid.UUID `json:"account_id"`
}
