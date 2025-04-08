package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/service/authorization"
	"github.com/google/uuid"
)

type MyServer struct {
	auth authorization.AuthorizationService
}

type Config struct {
	Port string
	DSN  string
}

var _ api.ServerInterface = (*MyServer)(nil)

func NewServer() *MyServer {
	auth := authorization.New()

	return &MyServer{
		auth: *auth,
	}
}

// Авторизация пользователя
// (POST /login)
func (s *MyServer) PostLogin(w http.ResponseWriter, r *http.Request) {

}

// Добавление товара в текущую приемку (только для сотрудников ПВЗ)
// (POST /products)
func (s *MyServer) PostProducts(w http.ResponseWriter, r *http.Request) {

}

// Получение списка ПВЗ с фильтрацией по дате приемки и пагинацией
// (GET /pvz)
func (s *MyServer) GetPvz(w http.ResponseWriter, r *http.Request, params api.GetPvzParams) {

}

// Создание ПВЗ (только для модераторов)
// (POST /pvz)
func (s *MyServer) PostPvz(w http.ResponseWriter, r *http.Request) {

}

// Закрытие последней открытой приемки товаров в рамках ПВЗ
// (POST /pvz/{pvzId}/close_last_reception)
func (s *MyServer) PostPvzPvzIdCloseLastReception(w http.ResponseWriter, r *http.Request, pvzId uuid.UUID) {

}

// Удаление последнего добавленного товара из текущей приемки (LIFO, только для сотрудников ПВЗ)
// (POST /pvz/{pvzId}/delete_last_product)
func (s *MyServer) PostPvzPvzIdDeleteLastProduct(w http.ResponseWriter, r *http.Request, pvzId uuid.UUID) {

}

// Создание новой приемки товаров (только для сотрудников ПВЗ)
// (POST /receptions)
func (s *MyServer) PostReceptions(w http.ResponseWriter, r *http.Request) {

}

// Регистрация пользователя
// (POST /register)
func (s *MyServer) PostRegister(w http.ResponseWriter, r *http.Request) {

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
