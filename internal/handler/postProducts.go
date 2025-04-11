package handler

import (
	"log/slog"
	"net/http"

	"github.com/Razzle131/pickupPoint/api"
)

// Добавление товара в текущую приемку (только для сотрудников ПВЗ)
// (POST /products)
func (s *MyServer) PostProducts(w http.ResponseWriter, r *http.Request) {
	slog.Debug("proccessing create product")
	defer slog.Debug("finished create product")

	token := r.Header.Get("Authorization")

	role, err := s.auth.ValidateToken(token)
	if err != nil {
		sendErrorResponse(w, "failed to validate token", http.StatusUnauthorized)
		return
	}

	if !validRole(string(role), api.UserRoleEmployee) {
		sendErrorResponse(w, "action forbidden", http.StatusForbidden)
		return
	}

	req, err := decodeBody[api.PostProductsJSONBody](r.Body)
	if err != nil {
		sendErrorResponse(w, "bad request body", http.StatusBadRequest)
		return
	}

	serviceRes, err := s.reception.AddProduct(string(req.Type), req.PvzId)
	if err != nil {
		sendErrorResponse(w, "failed to add product", http.StatusBadRequest)
		return
	}

	res := api.Product{
		DateTime:    &serviceRes.Date,
		Id:          &serviceRes.Id,
		ReceptionId: serviceRes.ReceptionId,
		Type:        serviceRes.Type,
	}

	sendInfoResponse(w, res)
}
