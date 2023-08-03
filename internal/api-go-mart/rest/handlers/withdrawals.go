package handlers

import (
	"errors"
	"net/http"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
)

func (h *Handler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromBody(w, r)

	withdrawals, err := h.db.Withdraw(userID)
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
	w.Header().Set("Content-type", "application/json")

	_, err = w.Write(withdrawals)
	if err != nil {
		h.log.Errorf("getOrders, cannot wrtie orders %s", string(withdrawals))
		return
	}
}
