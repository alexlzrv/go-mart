package loyalty

import (
	"encoding/json"
	"fmt"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/repository"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Accrual interface {
	GetActualInfo(order string) (*AccrualResponse, error)
}

type AccrualOrder struct {
	address string
	db      repository.Storage
	log     *zap.SugaredLogger
}

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func NewAccrual(address string, db repository.Storage, log *zap.SugaredLogger) *AccrualOrder {
	return &AccrualOrder{
		address: address,
		db:      db,
		log:     log,
	}
}

func (a *AccrualOrder) GetActualInfo(order string) (*AccrualResponse, error) {
	orderFromSystem, err := resty.New().SetRetryCount(5).R().Get(fmt.Sprintf("%s/api/orders/%s", a.address, order))
	if err != nil {
		a.log.Errorf("error while requesting for order %s: %s", order, err)
		return nil, err
	}
	var info AccrualResponse
	if err = json.Unmarshal(orderFromSystem.Body(), &info); err != nil {
		a.log.Errorf("error while unmarshalling order body: %s", err)
		return nil, err
	}
	return &info, nil
}
