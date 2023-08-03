package entities

import (
	"errors"
	"time"
)

var ErrNegativeBalance = errors.New("negative balance")

const (
	BalanceOperationRefill     string = "refill"
	BalanceOperationWithdrawal string = "withdrawal"
)

type Balance struct {
	UserID    int64   `json:"-"`
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type BalanceChange struct {
	UserID      int64     `json:"-"`
	Order       string    `json:"order"`
	Operation   string    `json:"-"`
	Amount      float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}
