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

// Авторизация пользователя
// (POST /login)
func (s *MyServer) PostLogin(w http.ResponseWriter, r *http.Request) {
	slog.Debug("processing login")
	defer slog.Debug("finished login")

	ctx, cancel := context.WithTimeout(context.Background(), consts.ContextTimeout)
	defer cancel()

	req, err := decodeBody[api.PostLoginJSONBody](r.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("body decode error: %s", err.Error()))
		sendErrorResponse(w, "bad creditionals", http.StatusUnauthorized)
		return
	}

	token, err := s.auth.Login(ctx, dto.UserCreditionalsDto{
		Email:    string(req.Email),
		Password: req.Password,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("login error: %s", err.Error()))
		sendErrorResponse(w, "bad creditionals", http.StatusUnauthorized)
		return
	}

	resp := api.Token(token)

	sendInfoResponse(w, resp, http.StatusOK)
}
