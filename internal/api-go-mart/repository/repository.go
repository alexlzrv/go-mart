package repository

import (
	"context"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
)

type Storage interface {
	Register(ctx context.Context, user *entities.User) error
	Login(user *entities.User) error

	GetUserOrders(userID int64) ([]byte, error)
	LoadOrder(order *entities.Order) error
	UpdateOrder(order *entities.Order) error
	GetNewAndProcessingOrder() ([]entities.Order, error)

	GetBalanceInfo(userID int64) ([]byte, error)
	Withdraw(userID int64) ([]byte, error)
	ChangeBalance(ctx context.Context, change *entities.BalanceChange) error
}
