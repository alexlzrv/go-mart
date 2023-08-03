package entities

import (
	"errors"
	"time"
)

const (
	OrderStatusNew        string = "NEW"
	OrderStatusProcessing string = "PROCESSING"
	OrderStatusInvalid    string = "INVALID"
	OrderStatusProcessed  string = "PROCESSED"
)

var (
	ErrOrderAlreadyAdded = errors.New("order has already been added")
	ErrOrderAddedByOther = errors.New("order has already been added by another user")
	ErrNoData            = errors.New("no data")
)

type Order struct {
	UserID     int64     `json:"-"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func NewOrder(userID int64, number string) *Order {
	return &Order{
		UserID: userID,
		Number: number,
		Status: OrderStatusNew,
	}
}
