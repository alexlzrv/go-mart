package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
	"github.com/alexlzrv/go-mart/internal/utils"
)

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	var user entities.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Errorf("error read body %s", err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Errorf("Cannot decode provided data: %s", err)
		return
	}

	if err = h.db.Register(r.Context(), &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Errorf("error with register user: %s", err)
		return
	}

	token, err := utils.GenerateToken(user.ID, h.key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Errorf("registration, error with generate token: %s", err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Header().Add("Authorization", token)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user entities.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Errorf("error read body %s", err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Errorf("Cannot decode provided data: %s", err)
		return
	}

	if err = h.db.Login(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Errorf("error with login user: %s", err)
		return
	}

	token, err := utils.GenerateToken(user.ID, h.key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Errorf("login, error with generate token: %s", err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Header().Add("Authorization", token)
	w.WriteHeader(http.StatusOK)
}
