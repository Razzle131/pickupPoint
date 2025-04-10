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
	"github.com/google/uuid"
)

type MyServer struct {
	auth      authorization.AuthorizationService
	pvz       pvz.PvzService
	reception reception.ReceptionService
}

type Config struct {
	Port string
	DSN  string
}

var _ api.ServerInterface = (*MyServer)(nil)

func NewServer(userRepo userRepo.UserRepo, pvzRepo pvzRepo.PvzRepo, productRepo productRepo.ProductRepo, receptionRepo receptionRepo.ReceptionRepo) *MyServer {
	pvzInfoRepo := pvzInfoRepo.NewCache(pvzRepo, productRepo, receptionRepo)

	auth := authorization.New(userRepo)
	pvz := pvz.New(pvzRepo, pvzInfoRepo)
	reception := reception.New(productRepo, receptionRepo)

	return &MyServer{
		auth:      *auth,
		pvz:       *pvz,
		reception: *reception,
	}
}

// Закрытие последней открытой приемки товаров в рамках ПВЗ
// (POST /pvz/{pvzId}/close_last_reception)
func (s *MyServer) PostPvzPvzIdCloseLastReception(w http.ResponseWriter, r *http.Request, pvzId uuid.UUID) {

}

// Удаление последнего добавленного товара из текущей приемки (LIFO, только для сотрудников ПВЗ)
// (POST /pvz/{pvzId}/delete_last_product)
func (s *MyServer) PostPvzPvzIdDeleteLastProduct(w http.ResponseWriter, r *http.Request, pvzId uuid.UUID) {

}

func sendErrorResponse(w http.ResponseWriter, errMsg string, status int) {
	resp, _ := json.Marshal(api.Error{Message: errMsg})
	slog.Error(errMsg)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func sendInfoResponse(w http.ResponseWriter, object any) {
	if object != nil {
		resp, err := json.Marshal(object)
		if err != nil {
			sendErrorResponse(w, "failed to form response", http.StatusInternalServerError)
			return
		}

		// to ensure that object is converted and there are no error and we dont have "superfluous response.WriteHeader call" message in log
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
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
