package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/repository/productRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzInfoRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/receptionRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/userRepo"
	"github.com/Razzle131/pickupPoint/internal/service/authorization"
	"github.com/Razzle131/pickupPoint/internal/service/pvz"
	"github.com/Razzle131/pickupPoint/internal/service/reception"
	"github.com/jmoiron/sqlx"
)

type MyServer struct {
	auth      authorization.AuthorizationService
	pvz       pvz.PvzService
	reception reception.ReceptionService
}

type Config struct {
	Port      string
	DbPort    string
	DbUser    string
	DbPasword string
	DbName    string
	DbHost    string
}

var _ api.ServerInterface = (*MyServer)(nil)

func NewServer(db *sqlx.DB) *MyServer {
	pvzRepo := pvzRepo.NewPg(db)
	receptionRepo := receptionRepo.NewPg(db)
	prodRepo := productRepo.NewPg(db)
	pvzInfoRepo := pvzInfoRepo.NewPg(db)
	userRepo := userRepo.NewPg(db)

	auth := authorization.New(userRepo)
	pvz := pvz.New(pvzRepo, pvzInfoRepo)
	reception := reception.New(prodRepo, receptionRepo, pvzRepo)

	return &MyServer{
		auth:      *auth,
		pvz:       *pvz,
		reception: *reception,
	}
}

func sendErrorResponse(w http.ResponseWriter, errMsg string, status int) {
	resp, _ := json.Marshal(api.Error{Message: errMsg})
	slog.Error(errMsg)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func sendInfoResponse(w http.ResponseWriter, object any, status int) {
	if object != nil {
		resp, err := json.Marshal(object)
		if err != nil {
			// 500 code is not presented in api schema, but i think giving it the badrequest status would be inappropriate
			sendErrorResponse(w, "failed to form response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	w.WriteHeader(status)
}

func validRole(having string, shouldBe api.UserRole) bool {
	return having == string(shouldBe)
}

func decodeBody[T any](body io.Reader) (T, error) {
	var res T

	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&res)
	if err != nil {
		return res, errors.New("bad request body")
	}

	slog.Debug("body decoded")
	return res, nil
}
