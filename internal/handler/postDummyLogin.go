package handler

import (
	"log/slog"
	"net/http"

	"github.com/Razzle131/pickupPoint/api"
)

// Получение тестового токена
// (POST /dummyLogin)
func (s *MyServer) PostDummyLogin(w http.ResponseWriter, r *http.Request) {
	slog.Debug("proccessing dummy auth")
	defer slog.Debug("finished dummy auth")

	req, err := decodeBody[api.PostDummyLoginJSONBody](r.Body)
	if err != nil {
		sendErrorResponse(w, "bad request body", http.StatusBadRequest)
		return
	}

	if !validRole(string(req.Role), api.UserRoleModerator) && !validRole(string(req.Role), api.UserRoleEmployee) {
		sendErrorResponse(w, "bad request body", http.StatusBadRequest)
		slog.Error("bad role")
		return
	}

	token, err := s.auth.DummyLogin(api.UserRole(req.Role))
	if err != nil {
		sendErrorResponse(w, "bad creditionals", http.StatusBadRequest)
		slog.Error("token creation error")
		return
	}

	resp := api.Token(token)

	sendInfoResponse(w, resp)
}
