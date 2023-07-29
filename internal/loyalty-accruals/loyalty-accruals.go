package loyalty

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
	"github.com/alexlzrv/go-mart/internal/api-go-mart/repository"
	"go.uber.org/zap"
)

type Accrual struct {
	address string
	db      repository.Storage
	log     *zap.SugaredLogger
}

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func New(address string, db repository.Storage, log *zap.SugaredLogger) *Accrual {
	return &Accrual{
		address: address,
		db:      db,
		log:     log,
	}
}

func (a *Accrual) updateOrdersInfo() error {
	requestContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := http.Client{}
	allOrder, err := a.db.GetAllOrder(requestContext)
	if err != nil {
		a.log.Errorf("error with get all orders %s", err)
		return err
	}

	for _, order := range allOrder {
		accrualOrder, err := a.sendOrder(requestContext, &client, order.Number)
		if err != nil {
			a.log.Errorf("sendOrder, error %s", err)
			return err
		}

		order.Status = entities.AccrualToOrderStatus(accrualOrder.Order)
		order.Accrual = accrualOrder.Accrual

		if err = a.db.UpdateOrder(requestContext, &order); err != nil {
			a.log.Errorf("error with update order info %s", err)
			return err
		}
	}

	return nil
}

func (a *Accrual) sendOrder(ctx context.Context, client *http.Client, order string) (*AccrualResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", a.address, order)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	accrualOrder := AccrualResponse{
		Accrual: 0,
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &accrualOrder)
	if err != nil {
		return nil, err
	}

	return &accrualOrder, nil
}
