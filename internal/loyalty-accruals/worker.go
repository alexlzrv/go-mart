package loyalty

import (
	"context"
	"time"
)

func (a *Accrual) OrderWorker(ctx context.Context) {
	a.log.Info("Start order worker")
	errorsCounter := 0
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			a.log.Info("Stop worker")
			return

		case <-ticker.C:
			if err := a.updateOrdersInfo(); err != nil {
				a.log.Errorf("updateOrdersInfo, error %s", err)
				errorsCounter++
				if errorsCounter > 10 {
					a.log.Infof("Stopping actualize orders because of many errors")
					return
				}
			}
		}
	}
}
