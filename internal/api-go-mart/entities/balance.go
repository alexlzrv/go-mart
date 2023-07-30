package entities

import (
	"errors"
	"time"
)

var ErrNegativeBalance = errors.New("negative balance")

type Balance struct {
	UserID    int64   `json:"-"`
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type BalanceChange struct {
	UserID      int64     `json:"-"`
	Order       string    `json:"order"`
	Amount      float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}
