package middleware

import (
	"context"
	"net/http"

	j "github.com/alexlzrv/go-mart/internal/utils"
	"go.uber.org/zap"
)

type Manager struct {
	key []byte
	log *zap.SugaredLogger
}

type key int

const (
	KeyPrincipalID key = iota
)

func NewManager(key []byte, log *zap.SugaredLogger) *Manager {
	return &Manager{
		key: key,
		log: log,
	}
}

func (mw *Manager) JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt := r.Header.Get("Authorization")
		if jwt == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := j.ParseToken(jwt, mw.key)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), KeyPrincipalID, user))

		next.ServeHTTP(w, r)
	})
}
