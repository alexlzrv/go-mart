package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
	j "github.com/alexlzrv/go-mart/internal/utils/jwt"
)

func (h *Handler) LoadOrders(w http.ResponseWriter, r *http.Request) {
	requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
	defer requestCancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("error read body %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")

	userID, err := j.ParseToken(authHeader)
	if err != nil {
		h.log.Errorf("error with parse token %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order := entities.NewOrder(userID, string(body))

	err = h.db.LoadOrder(requestContext, order)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidOrderNumber) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		if errors.Is(err, entities.ErrOrderAlreadyAdded) {
			w.WriteHeader(http.StatusOK)
			return
		}

		if errors.Is(err, entities.ErrOrderAddedByOther) {
			w.WriteHeader(http.StatusConflict)
			return
		}
	}
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
	defer requestCancel()

	authHeader := r.Header.Get("Authorization")

	userID, err := j.ParseToken(authHeader)
	if err != nil {
		h.log.Errorf("error with parse token %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orders, err := h.db.GetUserOrders(requestContext, userID)
	if err != nil {
		h.log.Errorf("getOrders, error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Infof("getOrders, order received user %d orders %s", userID, string(orders))

	w.Header().Set("content-type", "application/json")
	_, err = w.Write(orders)
	if err != nil {
		h.log.Errorf("getOrders, cannot wrtie orders %s", string(orders))
		return
	}
}
