package handlers

import (
	"net/http"
	"time"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/repository"
	j "github.com/alexlzrv/go-mart/internal/utils/jwt"
	"go.uber.org/zap"
)

const (
	requestTimeout = 5 * time.Second
)

type Handler struct {
	db  repository.Storage
	log *zap.SugaredLogger
}

func New(db repository.Storage, log *zap.SugaredLogger) *Handler {
	return &Handler{
		db:  db,
		log: log,
	}
}

func (h *Handler) getUserIDFromBody(w http.ResponseWriter, r *http.Request) (userID int64) {
	authHeader := r.Header.Get("Authorization")

	userID, err := j.ParseToken(authHeader)
	if err != nil {
		h.log.Errorf("error with parse token %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	return userID
}
