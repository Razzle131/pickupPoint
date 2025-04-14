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

// Удаление последнего добавленного товара из текущей приемки (LIFO, только для сотрудников ПВЗ)
// (POST /pvz/{pvzId}/delete_last_product)
func (s *MyServer) PostPvzPvzIdDeleteLastProduct(w http.ResponseWriter, r *http.Request, pvzId uuid.UUID) {
	slog.Debug("processing deleting product")
	defer slog.Debug("finished deleting product")

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

	_, err = s.reception.RemoveReceptionLastProduct(ctx, pvzId)
	if err != nil {
		slog.Error(fmt.Sprintf("remove last product error: %s", err.Error()))
		sendErrorResponse(w, "failed to remove product", http.StatusBadRequest)
		return
	}

	sendInfoResponse(w, nil, http.StatusOK)
}
