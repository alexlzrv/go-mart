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

	AccrualStatusRegistered string = "REGISTERED"
	AccrualStatusProcessing string = "PROCESSING"
	AccrualStatusInvalid    string = "INVALID"
	AccrualStatusProcessed  string = "PROCESSED"
)

var (
	ErrInvalidOrderNumber = errors.New("invalid order number")
	ErrOrderAlreadyAdded  = errors.New("order has already been added")
	ErrOrderAddedByOther  = errors.New("order has already been added by another user")
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

func AccrualToOrderStatus(status string) string {
	mapping := map[string]string{
		AccrualStatusRegistered: OrderStatusNew,
		AccrualStatusProcessing: OrderStatusProcessing,
		AccrualStatusInvalid:    OrderStatusInvalid,
		AccrualStatusProcessed:  OrderStatusProcessed,
	}

	return mapping[status]
}
