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

	GetBalanceInfo(ctx context.Context, userID int64) ([]byte, error)
	Withdraw(ctx context.Context, userID int64) ([]byte, error)
	GetWithdrawals(ctx context.Context, change *entities.BalanceChange) error
}
