package apigomart

import (
	"net/http"
	"sync"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/repository/pgrepo"
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/routers"
	"github.com/alexlzrv/go-mart/internal/config"
	"github.com/alexlzrv/go-mart/internal/logger"
	"github.com/alexlzrv/go-mart/sql"
	"github.com/go-chi/chi/v5"
)

func Run() {
	cfg := config.NewConfig()

	log, err := logger.LogInitializer(cfg.LogLevel)
	if err != nil {
		return
	}

	log.Infof("Run server at address %s", cfg.ServerAddress)
	var (
		r   = chi.NewRouter()
		srv = &http.Server{
			Addr:    cfg.ServerAddress,
			Handler: r,
		}
	)

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

	routers.GetRoutes(r, repo, log)

	log.Info("Server is running...")
	if err = srv.ListenAndServe(); err != nil {
		log.Fatalf("Error with server running: %v", err)
	}

	wg := &sync.WaitGroup{}

	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	log.Info("Server is running...")
	//	if err = srv.ListenAndServe(); err != nil {
	//		log.Fatalf("Error with server running: %v", err)
	//	}
	//
	//}()
	//
	//if err = srv.Shutdown(context.Background()); err != nil {
	//	log.Errorf("server shutdown %v", err)
	//	return
	//}

	wg.Wait()
}
