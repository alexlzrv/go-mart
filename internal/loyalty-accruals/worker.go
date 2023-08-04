package loyalty

import (
	"context"
	"errors"
	"time"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
	"github.com/alexlzrv/go-mart/internal/api-go-mart/repository"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Worker struct {
	db      repository.Storage
	accrual Accrual
	log     *zap.SugaredLogger
	workers int
}

func NewWorker(db repository.Storage, accrual Accrual, workersCount int, log *zap.SugaredLogger) *Worker {
	return &Worker{
		db:      db,
		workers: workersCount,
		accrual: accrual,
		log:     log,
	}
}

func (w *Worker) Run(ctx context.Context) {
	w.log.Info("Start order worker")
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		orders, err := w.db.GetNewAndProcessingOrder()
		if err != nil {
			w.log.Errorf("error with get processing orders %s", err)
		}

		task := w.generateOrderTask(ctx, orders)

		gr, grCtx := errgroup.WithContext(context.Background())

		for i := 0; i < w.workers; i++ {
			gr.Go(func() error {
				return w.orderTaskWorker(grCtx, task)
			})
		}

		if err = gr.Wait(); err != nil {
			var accrualErr *AccrualCoolDownError
			if errors.As(err, &accrualErr) {
				time.Sleep(accrualErr.CoolDown)
				continue
			} else {
				w.log.Errorf("workers run with error %s", err)
			}
		}
	}
}

func (w *Worker) generateOrderTask(ctx context.Context, orders []entities.Order) chan entities.Order {
	ordersChan := make(chan entities.Order, w.workers)

	go func() {
		defer close(ordersChan)

		for i := range orders {
			select {
			case <-ctx.Done():
				return
			case ordersChan <- orders[i]:
			}
		}
	}()

	return ordersChan
}

func (w *Worker) orderTaskWorker(ctx context.Context, task <-chan entities.Order) error {
	for order := range task {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := w.db.UpdateOrder(&entities.Order{Status: entities.OrderStatusProcessing, Number: order.Number})
			if err != nil {
				w.log.Errorf("errror with update orders info %s", err)
				return err
			}

			orderInfo, err := w.accrual.GetActualInfo(order.Number)
			if err != nil {
				if errors.Is(err, ErrAccrualOrderNotFound) {
					continue
				}
				return err
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
					return err
				}

				err = w.db.ChangeBalance(ctx, &entities.BalanceChange{
					UserID:    order.UserID,
					Order:     orderInfo.Order,
					Operation: entities.BalanceOperationRefill,
					Amount:    orderInfo.Accrual,
				})

				if err != nil {
					w.log.Errorf("errror with get withdrawals %s", err)
					return err
				}
			}
		}
	}

	return nil
}
