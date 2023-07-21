package orders

import (
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/handlers"
	"github.com/go-chi/chi/v5"
)

func Orders(r chi.Router, handler *handlers.Handler) {
	r.Group(func(router chi.Router) {
		router.Post("/orders", handler.LoadOrders)
		router.Get("/orders", handler.GetOrders)
	})
}
