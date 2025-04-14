package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/google/uuid"
)

// Закрытие последней открытой приемки товаров в рамках ПВЗ
// (POST /pvz/{pvzId}/close_last_reception)
func (s *MyServer) PostPvzPvzIdCloseLastReception(w http.ResponseWriter, r *http.Request, pvzId uuid.UUID) {
	slog.Debug("processing closing reception")
	defer slog.Debug("finished closing reception")

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

	serviceRes, err := s.reception.CloseReception(ctx, pvzId)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to close reception error: %s", err.Error()))
		sendErrorResponse(w, "failed to close reception", http.StatusBadRequest)
		return
	}

	res := serviceRes.ToApiModel()

	sendInfoResponse(w, res, http.StatusOK)
}
