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

// Добавление товара в текущую приемку (только для сотрудников ПВЗ)
// (POST /products)
func (s *MyServer) PostProducts(w http.ResponseWriter, r *http.Request) {
	slog.Debug("processing create product")
	defer slog.Debug("finished create product")

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

	req, err := decodeBody[api.PostProductsJSONBody](r.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("decode body error: %s", err.Error()))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	product := dto.ProductDto{
		Type: string(req.Type),
	}

	if !product.IsValidType() {
		slog.Error(fmt.Sprintf("bad product type: %s", product.Type))
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	serviceRes, err := s.reception.AddProduct(ctx, product.Type, req.PvzId)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to add product error: %s", err.Error()))
		sendErrorResponse(w, "failed to add product", http.StatusBadRequest)
		return
	}

	res := serviceRes.ToApiModel()

	sendInfoResponse(w, res, http.StatusCreated)
}
