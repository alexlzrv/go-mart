package middleware

import (
	"net/http"

	j "github.com/alexlzrv/go-mart/internal/utils"
	"github.com/golang-jwt/jwt/v4"
)

type Manager struct {
	key []byte
}

func NewManager(key []byte) *Manager {
	return &Manager{
		key: key,
	}
}

func (mw *Manager) JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := &j.Claims{}

		token, err := jwt.ParseWithClaims(tokenHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return mw.key, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Add("Authorization", tokenHeader)

		next.ServeHTTP(w, r)
	})
}
