package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/middleware"
	"github.com/alexlzrv/go-mart/internal/utils"
)

func (h *Handler) LoadOrders(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("error read body %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middleware.KeyPrincipalID).(int64)

	order := entities.NewOrder(userID, string(body))

	if ok := utils.LuhnCheck(order.Number); !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = h.db.LoadOrder(order)
	if err != nil {
		if errors.Is(err, entities.ErrOrderAlreadyAdded) {
			w.WriteHeader(http.StatusOK)
			return
		}

		if errors.Is(err, entities.ErrOrderAddedByOther) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		h.log.Errorf("error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.KeyPrincipalID).(int64)

	orders, err := h.db.GetUserOrders(userID)
	if err != nil {
		if errors.Is(err, entities.ErrNoData) {
			h.log.Errorf("getOrders, no data")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.log.Errorf("getOrders, error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Infof("getOrders, order received user %d orders %s", userID, string(orders))

	w.Header().Set("Content-type", "application/json")
	_, err = w.Write(orders)
	if err != nil {
		h.log.Errorf("getOrders, cannot wrtie orders %s", string(orders))
		return
	}
}
