package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
)

func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
	defer requestCancel()

	userID := h.getUserIDFromBody(w, r)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("error read body %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var change entities.BalanceChange

	err = json.Unmarshal(body, &change)
	if err != nil {
		h.log.Errorf("getWithdrawals, error with marshal %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	change.UserID = userID

	err = h.db.GetWithdrawals(requestContext, &change)
	if err != nil {
		if errors.Is(err, entities.ErrNegativeBalance) {
			h.log.Errorf("negative balance")
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		if errors.Is(err, entities.ErrInvalidOrderNumber) {
			h.log.Errorf("invalid order number")
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		h.log.Errorf("getWithdrawals, error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
	defer requestCancel()

	userID := h.getUserIDFromBody(w, r)

	balance, err := h.db.GetBalanceInfo(requestContext, userID)
	if err != nil {
		h.log.Errorf("getBalanceInfo, error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	_, err = w.Write(balance)
	if err != nil {
		h.log.Errorf("getBalanceInfo, cannot wrtie orders %s", string(balance))
		return
	}
}
