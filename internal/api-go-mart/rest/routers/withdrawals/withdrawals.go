package withdrawals

import (
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/handlers"
	"github.com/go-chi/chi/v5"
)

func Withdrawals(r chi.Router, handler *handlers.Handler) {
	r.Group(func(router chi.Router) {
		router.Get("/withdrawals", handler.Withdrawals)
	})
}
