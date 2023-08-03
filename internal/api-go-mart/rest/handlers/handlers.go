package handlers

import (
	"github.com/alexlzrv/go-mart/internal/api-go-mart/repository"
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
