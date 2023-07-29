package apigomart

import (
	"context"
	"net/http"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/repository/pgrepo"
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/routers"
	"github.com/alexlzrv/go-mart/internal/config"
	"github.com/alexlzrv/go-mart/internal/logger"
	"github.com/alexlzrv/go-mart/internal/loyalty-accruals"
	"github.com/alexlzrv/go-mart/internal/runner"
	"github.com/alexlzrv/go-mart/sql"
)

func Run() {
	ctx := context.Background()
	cfg := config.NewConfig()

	log, err := logger.LogInitializer(cfg.LogLevel)
	if err != nil {
		return
	}

	pg, err := sql.NewPostgresStorage(cfg.DSN, log)
	if err != nil {
		log.Errorf("Error init postgres storage %s", err)
	}

	log.Infof("Database connection open")

	defer func() {
		err = pg.Close()
		if err != nil {
			return
		}
		log.Infof("Database connection closed")
	}()

	repo := pgrepo.NewRepository(pg, log)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: routers.NewRoutes(repo, log),
	}

	loyaltyAccrual := loyalty.New(cfg.AccrualAddress, repo, log)

	r := runner.New(server, loyaltyAccrual, log)
	if err = r.Run(ctx); err != nil {
		log.Errorf("error while running runner: %s", err)
		return
	}
}
