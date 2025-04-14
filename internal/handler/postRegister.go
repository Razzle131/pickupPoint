package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/dto"
)

// Регистрация пользователя
// (POST /register)
func (s *MyServer) PostRegister(w http.ResponseWriter, r *http.Request) {
	slog.Debug("processing register")
	defer slog.Debug("finished register")

	ctx, cancel := context.WithTimeout(context.Background(), consts.ContextTimeout)
	defer cancel()

	req, err := decodeBody[api.PostRegisterJSONBody](r.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("body decode error: %s", err.Error()))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	userDto := dto.UserCreditionalsDto{
		Email:    string(req.Email),
		Password: req.Password,
		Role:     string(req.Role),
	}

	if !userDto.IsValidRole() {
		slog.Error(fmt.Sprintf("bad role: %s", req.Role))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	user, err := s.auth.Register(ctx, userDto)
	if err != nil {
		slog.Error(fmt.Sprintf("register error: %s", err.Error()))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	res := user.ToApiModel()

	sendInfoResponse(w, res, http.StatusCreated)
}
