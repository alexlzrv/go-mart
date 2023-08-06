package repository

import (
	"context"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
)

type Storage interface {
	Register(ctx context.Context, user *entities.User) error
	Login(ctx context.Context, user *entities.User) error

	GetUserOrders(ctx context.Context, userID int64) ([]byte, error)
	LoadOrder(ctx context.Context, order *entities.Order) error
	UpdateOrder(ctx context.Context, order *entities.Order) error
	GetAllOrder(ctx context.Context) ([]entities.Order, error)

	GetBalanceInfo(login string) ([]byte, error)
	Withdraw(login string, orderID string, sum float64) error
	GetWithdrawals(login string) ([]byte, error)
}
