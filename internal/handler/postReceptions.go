package handler

import (
	"log/slog"
	"net/http"

	"github.com/Razzle131/pickupPoint/api"
)

// Создание новой приемки товаров (только для сотрудников ПВЗ)
// (POST /receptions)
func (s *MyServer) PostReceptions(w http.ResponseWriter, r *http.Request) {
	slog.Debug("proccessing create reception")
	defer slog.Debug("finished create reception")

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

	req, err := decodeBody[api.PostReceptionsJSONBody](r.Body)
	if err != nil {
		sendErrorResponse(w, "bad request body", http.StatusBadRequest)
		return
	}

	serviceRes, err := s.reception.AddReception(req.PvzId)
	if err != nil {
		sendErrorResponse(w, "failed to add reception", http.StatusBadRequest)
		return
	}

	res := api.Reception{
		DateTime: serviceRes.Date,
		Id:       &serviceRes.Id,
		PvzId:    serviceRes.PvzId,
		Status:   serviceRes.Status,
	}

	sendInfoResponse(w, res)
}
