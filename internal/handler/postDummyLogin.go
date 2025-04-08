package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Razzle131/pickupPoint/api"
)

// Получение тестового токена
// (POST /dummyLogin)
func (s *MyServer) PostDummyLogin(w http.ResponseWriter, r *http.Request) {
	slog.Debug("proccessing dummy auth")
	defer slog.Debug("finished dummy auth")

	var authStruct api.PostDummyLoginJSONBody

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&authStruct)
	if err != nil {
		sendErrorResponse(w, "bad request body", http.StatusBadRequest)
		slog.Error("bad request body")
		return
	}
	slog.Debug("body decoded")

	// role validation
	// if (authStruct.Role != api.PostDummyLoginJSONBodyRole(api.UserRoleModerator)) || (api.UserRole(authStruct.Role) != api.UserRoleEmployee) {
	// 	slog.Error("bad request body")
	// 	return
	// }

	token, err := s.auth.DummyLogin(api.UserRole(authStruct.Role))
	if err != nil {
		sendErrorResponse(w, "bad creditionals", http.StatusBadRequest)
		slog.Error("token creation error")
		return
	}

	resp := api.Token(token)

	sendInfoResponse(w, resp)
}
