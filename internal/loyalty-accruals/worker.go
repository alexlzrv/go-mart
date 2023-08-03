package loyalty

import (
	"context"
	"time"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
	"github.com/alexlzrv/go-mart/internal/api-go-mart/repository"
	"go.uber.org/zap"
)

type Worker struct {
	ordersChan chan *entities.Order
	db         repository.Storage
	accrual    Accrual
	log        *zap.SugaredLogger
}

const lenChan = 10

func NewWorker(db repository.Storage, accrual Accrual, log *zap.SugaredLogger) *Worker {
	ordersChan := make(chan *entities.Order, lenChan)

	w := &Worker{
		ordersChan: ordersChan,
		db:         db,
		accrual:    accrual,
		log:        log,
	}

	go w.orderWorker()

	return w
}

func (w *Worker) orderWorker() {
	w.log.Info("Start order worker")
	ticker := time.NewTicker(lenChan * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		orders, err := w.db.GetNewAndProcessingOrder()
		if err != nil {
			w.log.Errorf("error with get processing orders %s", err)
		}

		for i := range orders {
			w.ordersChan <- &orders[i]
		}
	}
}

func (w *Worker) Worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			w.log.Info("Stopping order worker")
			return
		case order := <-w.ordersChan:
			err := w.db.UpdateOrder(&entities.Order{Status: entities.OrderStatusProcessing, Number: order.Number})
			if err != nil {
				w.log.Errorf("errror with update orders info %s", err)
				return
			}

			orderInfo, err := w.accrual.GetActualInfo(order.Number)

			//проверить тип ошибки ретрай афтер

			//errGroup
			if err != nil {
				w.log.Errorf("error with get actual info %s", err)
				return
			}

			if orderInfo.Status == entities.OrderStatusProcessed {
				err = w.db.UpdateOrder(&entities.Order{
					Number:  orderInfo.Order,
					UserID:  order.UserID,
					Accrual: orderInfo.Accrual,
					Status:  orderInfo.Status,
				})

				if err != nil {
					w.log.Errorf("errror with update orders infoo %s", err)
					return
				}

				err = w.db.ChangeBalance(ctx, &entities.BalanceChange{
					UserID:    order.UserID,
					Order:     orderInfo.Order,
					Operation: entities.BalanceOperationRefill,
					Amount:    orderInfo.Accrual,
				})

				if err != nil {
					w.log.Errorf("errror with get withdrawals %s", err)
					return
				}
			}
		}
	}
}
