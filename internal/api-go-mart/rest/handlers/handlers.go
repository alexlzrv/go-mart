package handlers

import (
	"time"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/repository"
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
