package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
	"github.com/alexlzrv/go-mart/internal/api-go-mart/rest/middleware"
	"github.com/alexlzrv/go-mart/internal/utils"
)

func (h *Handler) ChangeBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.KeyPrincipalID).(int64)

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
	h.log.Infof("changeBalance request body %v", change)

	change.UserID = userID
	change.Operation = entities.BalanceOperationWithdrawal

	if ok := utils.LuhnCheck(change.Order); !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = h.db.ChangeBalance(r.Context(), &change)
	if err != nil {
		if errors.Is(err, entities.ErrNegativeBalance) {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		h.log.Errorf("getWithdrawals, error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.KeyPrincipalID).(int64)

	balance, err := h.db.GetBalanceInfo(userID)
	if err != nil {
		h.log.Errorf("getBalanceInfo, error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Infof("getBalance request body %s", string(balance))

	w.Header().Set("Content-type", "application/json")
	_, err = w.Write(balance)
	if err != nil {
		h.log.Errorf("getBalanceInfo, cannot wrtie orders %s", string(balance))
		return
	}
}
