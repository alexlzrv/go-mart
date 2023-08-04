package loyalty

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/repository"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Accrual interface {
	GetActualInfo(order string) (*AccrualResponse, error)
}

type AccrualCoolDownError struct {
	CoolDown time.Duration
}

func (e *AccrualCoolDownError) Error() string {
	return fmt.Sprintf("wait %v", e.CoolDown.Seconds())
}

var (
	ErrAccrualOrderNotFound = errors.New("accrual order not found")
)

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

func newAccrualCoolDownError(coolDown time.Duration) *AccrualCoolDownError {
	return &AccrualCoolDownError{
		CoolDown: coolDown,
	}
}

func NewAccrual(address string, db repository.Storage, log *zap.SugaredLogger) *AccrualOrder {
	return &AccrualOrder{
		address: address,
		db:      db,
		log:     log,
	}
}

func (a *AccrualOrder) GetActualInfo(order string) (*AccrualResponse, error) {
	resp, err := resty.New().SetRetryCount(5).R().Get(fmt.Sprintf("%s/api/orders/%s", a.address, order))
	if err != nil {
		a.log.Errorf("error while requesting for order %s: %s", order, err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == http.StatusNoContent {
			return nil, ErrAccrualOrderNotFound
		}

		if resp.StatusCode() == http.StatusTooManyRequests {
			coolDown, err := strconv.Atoi(resp.Header().Get("Retry-After"))
			if err != nil {
				return nil, err
			}

			return nil, newAccrualCoolDownError(time.Duration(coolDown) * time.Second)
		}
	}

	var info AccrualResponse
	if err = json.Unmarshal(resp.Body(), &info); err != nil {
		a.log.Errorf("error while unmarshalling order body: %s", err)
		return nil, err
	}
	return &info, nil
}
