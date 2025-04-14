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

// Создание ПВЗ (только для модераторов)
// (POST /pvz)
func (s *MyServer) PostPvz(w http.ResponseWriter, r *http.Request) {
	slog.Debug("processing create pvz")
	defer slog.Debug("finished create pvz")

	ctx, cancel := context.WithTimeout(context.Background(), consts.ContextTimeout)
	defer cancel()

	token := r.Header.Get("Authorization")

	role, err := s.auth.ValidateToken(ctx, token)
	if err != nil {
		slog.Error(fmt.Sprintf("validate token error: %s", err.Error()))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	if !validRole(string(role), api.UserRoleModerator) {
		slog.Error(fmt.Sprintf("bad role: %s", string(role)))
		sendErrorResponse(w, "action forbidden", http.StatusForbidden)
		return
	}

	req, err := decodeBody[api.PostPvzJSONRequestBody](r.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("body decode error: %s", err.Error()))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	dto := dto.PvzDto{}
	dto.FromApiModel(req)

	if !dto.IsValidCity() {
		slog.Error(fmt.Sprintf("bad request city: %s", req.City))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	serviceRes, err := s.pvz.CreatePvz(ctx, dto)
	if err != nil {
		slog.Error(fmt.Sprintf("create pvz error: %s", err.Error()))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	res := serviceRes.ToApiModel()

	sendInfoResponse(w, res, http.StatusCreated)
}
