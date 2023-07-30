package runner

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexlzrv/go-mart/internal/config"
	"github.com/alexlzrv/go-mart/internal/loyalty-accruals"
	"go.uber.org/zap"
)

type Runner struct {
	log            *zap.SugaredLogger
	server         *http.Server
	loyaltyAccrual *loyalty.Accrual
	cfq            *config.Config
}

func New(server *http.Server, loyaltyAccrual *loyalty.Accrual, log *zap.SugaredLogger, cfg *config.Config) *Runner {
	return &Runner{
		server:         server,
		log:            log,
		loyaltyAccrual: loyaltyAccrual,
		cfq:            cfg,
	}
}

func (r *Runner) Run(ctx context.Context) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		r.log.Infof("Stopping server")
		if err := r.server.Shutdown(ctx); err != nil {
			r.log.Errorf("Error stopping server: %s", err)
		}
	}()

	if r.cfq.AccrualAddress != "" {
		go r.loyaltyAccrual.OrderWorker(ctx)
	}

	r.log.Infof("Starting server on addr: %s", r.server.Addr)
	r.log.Infof("Starting accrual server on addr: %s", r.cfq.AccrualAddress)
	if err := r.server.ListenAndServe(); err != nil {
		r.log.Errorf("error while running server: %s", err)
		return err
	}
	return nil
}
