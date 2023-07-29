package routers

import (
	"github.com/alexlzrv/go-mart/internal/api-go-mart/repository"
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/handlers"
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/middleware"
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/routers/auth"
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/routers/balance"
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/routers/orders"
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/routers/withdrawals"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRoutes(db repository.Storage, log *zap.SugaredLogger) *chi.Mux {
	handler := handlers.New(db, log)
	r := chi.NewRouter()
	r.Route("/api/user/", func(r chi.Router) {
		auth.AuthRouters(r, handler)

		r.Group(func(router chi.Router) {
			router.Use(middleware.JWTAuth)
			orders.Orders(router, handler)
			balance.Balance(router, handler)
			withdrawals.Withdrawals(router, handler)
		})
	})

	return r
}
