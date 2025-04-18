package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/consts"
)

// Создание новой приемки товаров (только для сотрудников ПВЗ)
// (POST /receptions)
func (s *MyServer) PostReceptions(w http.ResponseWriter, r *http.Request) {
	slog.Debug("processing create reception")
	defer slog.Debug("finished create reception")

	ctx, cancel := context.WithTimeout(context.Background(), consts.ContextTimeout)
	defer cancel()

	token := r.Header.Get("Authorization")

	role, err := s.auth.ValidateToken(ctx, token)
	if err != nil {
		slog.Error(fmt.Sprintf("validate token error: %s", err.Error()))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	if !validRole(string(role), api.UserRoleEmployee) {
		slog.Error(fmt.Sprintf("bad role: %s", string(role)))
		sendErrorResponse(w, "action forbidden", http.StatusForbidden)
		return
	}

	req, err := decodeBody[api.PostReceptionsJSONBody](r.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("decode body error: %s", err.Error()))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	serviceRes, err := s.reception.AddReception(ctx, req.PvzId)
	if err != nil {
		slog.Error(fmt.Sprintf("add reception error: %s", err.Error()))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	res := serviceRes.ToApiModel()

	sendInfoResponse(w, res, http.StatusCreated)
}
