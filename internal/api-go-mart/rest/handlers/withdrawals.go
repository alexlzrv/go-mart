package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
)

func (h *Handler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
	defer requestCancel()

	userID := h.getUserIDFromBody(w, r)

	withdrawals, err := h.db.Withdraw(requestContext, userID)
	if err != nil {
		if errors.Is(err, entities.ErrNoData) {
			h.log.Errorf("withdrawals, no data")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.log.Errorf("withdrawals, error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")

	_, err = w.Write(withdrawals)
	if err != nil {
		h.log.Errorf("getOrders, cannot wrtie orders %s", string(withdrawals))
		return
	}
}
