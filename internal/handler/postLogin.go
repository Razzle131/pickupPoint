package handler

import (
	"log/slog"
	"net/http"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/dto"
)

// Авторизация пользователя
// (POST /login)
func (s *MyServer) PostLogin(w http.ResponseWriter, r *http.Request) {
	slog.Debug("proccessing login")
	defer slog.Debug("finished login")

	req, err := decodeBody[api.PostLoginJSONBody](r.Body)
	if err != nil {
		sendErrorResponse(w, "bad request body", http.StatusBadRequest)
		return
	}

	token, err := s.auth.Login(dto.LoginDto{
		Email:    string(req.Email),
		Password: req.Password,
	})
	if err != nil {
		sendErrorResponse(w, "bad creditionals", http.StatusBadRequest)
		slog.Error("token creation error")
		return
	}

	resp := api.Token(token)

	sendInfoResponse(w, resp)
}
