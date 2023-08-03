package handlers

import (
	"net/http"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/repository"
	j "github.com/alexlzrv/go-mart/internal/utils"
	"go.uber.org/zap"
)

type Handler struct {
	db  repository.Storage
	log *zap.SugaredLogger
	key []byte
}

func NewHandler(db repository.Storage, log *zap.SugaredLogger, key []byte) *Handler {
	return &Handler{
		db:  db,
		log: log,
		key: key,
	}
}

func (h *Handler) getUserIDFromBody(w http.ResponseWriter, r *http.Request) (userID int64) {
	authHeader := r.Header.Get("Authorization")

	userID, err := j.ParseToken(authHeader, h.key)
	if err != nil {
		h.log.Errorf("error with parse token %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	return userID
}
