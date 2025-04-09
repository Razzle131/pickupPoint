package handler

import (
	"log/slog"
	"net/http"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/oapi-codegen/runtime/types"
)

// Регистрация пользователя
// (POST /register)
func (s *MyServer) PostRegister(w http.ResponseWriter, r *http.Request) {
	slog.Debug("proccessing register")
	defer slog.Debug("finished register")

	req, err := decodeBody[api.PostRegisterJSONBody](r.Body)
	if err != nil {
		sendErrorResponse(w, "bad request body", http.StatusBadRequest)
		return
	}

	if !validRole(string(req.Role), api.UserRoleModerator) && !validRole(string(req.Role), api.UserRoleEmployee) {
		sendErrorResponse(w, "bad request body", http.StatusBadRequest)
		slog.Error("bad role")
		return
	}

	user, err := s.auth.Register(dto.UserDto{
		Email:    string(req.Email),
		Password: req.Password,
		Role:     api.UserRole(req.Role),
	})
	if err != nil {
		sendErrorResponse(w, "failed to register user", http.StatusBadRequest)
		return
	}

	id := types.UUID(user.Id)
	res := api.User{
		Email: types.Email(user.Email),
		Id:    &id,
		Role:  user.Role,
	}

	sendInfoResponse(w, res)
}
