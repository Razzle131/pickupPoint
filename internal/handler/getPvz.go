package handler

import (
	"net/http"

	"github.com/Razzle131/pickupPoint/api"
)

// Получение списка ПВЗ с фильтрацией по дате приемки и пагинацией
// (GET /pvz)
func (s *MyServer) GetPvz(w http.ResponseWriter, r *http.Request, params api.GetPvzParams) {

}
