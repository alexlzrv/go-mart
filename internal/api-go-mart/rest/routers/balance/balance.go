package balance

import (
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/handlers"
	"github.com/go-chi/chi/v5"
)

func Balance(r chi.Router, handler *handlers.Handler) {
	r.Group(func(router chi.Router) {
		router.Post("/balance/withdraw", handler.BalanceWithdrawals)
		router.Get("/balance", handler.GetBalance)
	})
}
