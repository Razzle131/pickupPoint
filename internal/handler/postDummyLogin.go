package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/consts"
)

// Получение тестового токена
// (POST /dummyLogin)
func (s *MyServer) PostDummyLogin(w http.ResponseWriter, r *http.Request) {
	slog.Debug("processing dummy auth")
	defer slog.Debug("finished dummy auth")

	ctx, cancel := context.WithTimeout(context.Background(), consts.ContextTimeout)
	defer cancel()

	req, err := decodeBody[api.PostDummyLoginJSONBody](r.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("body decode error: %s", err.Error()))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	if !validRole(string(req.Role), api.UserRoleModerator) && !validRole(string(req.Role), api.UserRoleEmployee) {
		slog.Error(fmt.Sprintf("bad role: %s", req.Role))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	token, err := s.auth.DummyLogin(ctx, string(req.Role))
	if err != nil {
		slog.Error(fmt.Sprintf("token creation error: %s", err.Error()))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	resp := api.Token(token)

	sendInfoResponse(w, resp, http.StatusOK)
}
