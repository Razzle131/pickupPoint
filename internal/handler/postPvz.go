package handler

import (
	"log/slog"
	"net/http"

	"github.com/Razzle131/pickupPoint/api"
)

// Создание ПВЗ (только для модераторов)
// (POST /pvz)
func (s *MyServer) PostPvz(w http.ResponseWriter, r *http.Request) {
	slog.Debug("proccessing create pvz")
	defer slog.Debug("finished create pvz")

	token := r.Header.Get("Authorization")

	role, err := s.auth.ValidateToken(token)

	slog.Debug(string(role))

	if err != nil {
		sendErrorResponse(w, "bad token", http.StatusBadRequest)
		return
	}

	if !validRole(string(role), api.UserRoleModerator) {
		sendErrorResponse(w, "bad role", http.StatusForbidden)
		return
	}

	req, err := decodeBody[api.PostPvzJSONRequestBody](r.Body)
	if err != nil {
		sendErrorResponse(w, "bad request body", http.StatusBadRequest)
		return
	}

	res, err := s.pvz.CreatePvz(req)
	if err != nil {
		sendErrorResponse(w, "bad request body", http.StatusBadRequest)
		return
	}

	sendInfoResponse(w, res)
}
