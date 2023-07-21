package auth

import (
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/handlers"
	"github.com/go-chi/chi/v5"
)

func AuthRouters(r chi.Router, handler *handlers.Handler) {
	r.Group(func(router chi.Router) {
		router.Post("/register", handler.Registration)
		router.Post("/login", handler.Login)
	})
}
