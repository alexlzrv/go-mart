package apigomart

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/repository/pgrepo"
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/routers"
	"github.com/alexlzrv/go-mart/internal/config"
	"github.com/alexlzrv/go-mart/internal/logger"
	"github.com/alexlzrv/go-mart/internal/loyalty-accruals"
	"github.com/alexlzrv/go-mart/sql"
)

func Run(cfg *config.Config) {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	log, err := logger.LogInitializer(cfg.LogLevel)
	if err != nil {
		return
	}

	pg, err := sql.NewStorage(cfg.DSN, log)
	if err != nil {
		log.Errorf("Error init postgres storage %s", err)
		return
	}

	log.Infof("Database connection open")

	repo := pgrepo.NewRepository(pg, log)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: routers.NewRoutes(repo, log, []byte(cfg.SecretKey)),
	}

	wg := &sync.WaitGroup{}
	defer func() {
		wg.Wait()
	}()

	wg.Add(1)
	go func() {
		defer log.Info("Database connection closed")
		defer wg.Done()
		<-ctx.Done()

		pg.Close()
	}()

	loyaltyAccrual := loyalty.NewAccrual(cfg.AccrualAddress, repo, log)
	worker := loyalty.NewWorker(repo, loyaltyAccrual, log)

	wg.Add(1)
	go func() {
		worker.Worker(ctx)
	}()

	componentsErrs := make(chan error, 1)

	go func(errs chan<- error) {
		log.Infof("Starting server on addr: %s", server.Addr)
		log.Infof("Starting accrual server on addr: %s", cfg.AccrualAddress)
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			errs <- fmt.Errorf("listen and server has failed: %w", err)
		}
	}(componentsErrs)

	wg.Add(1)
	go func() {
		defer log.Infof("Stopping server")
		defer wg.Done()
		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdownTimeoutCtx()
		if err := server.Shutdown(shutdownTimeoutCtx); err != nil {
			log.Infof("an error occurred during server shutdown: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-componentsErrs:
		log.Info(err)
		cancelCtx()
	}

	go func() {
		ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelCtx()

		<-ctx.Done()
		log.Fatal("failed to gracefully shutdown the service")
	}()
}
