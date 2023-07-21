package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
	"github.com/alexlzrv/go-mart/internal/utils/jwt"
)

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
	defer requestCancel()

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

	if err = h.db.Register(requestContext, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Errorf("error with register user: %s", err)
		return
	}

	if err = h.db.Login(requestContext, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Errorf("error with login user after registration: %s", err)
		return
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Errorf("registration, error with generate token: %s", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Header().Add("Authorization", token)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
	defer requestCancel()

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

	if err = h.db.Login(requestContext, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Errorf("error with login user: %s", err)
		return
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Errorf("login, error with generate token: %s", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Header().Add("Authorization", token)
	w.WriteHeader(http.StatusOK)
}
