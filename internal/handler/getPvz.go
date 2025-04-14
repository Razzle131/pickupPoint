package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/dto"
)

// There are no errors in api schema, so I wont send any respond if the error occurres

// Получение списка ПВЗ с фильтрацией по дате приемки и пагинацией
// (GET /pvz)
func (s *MyServer) GetPvz(w http.ResponseWriter, r *http.Request, params api.GetPvzParams) {
	slog.Debug("processing getting pvz")
	defer slog.Debug("finished getting pvz")

	ctx, cancel := context.WithTimeout(context.Background(), consts.ContextTimeout)
	defer cancel()

	token := r.Header.Get("Authorization")

	role, err := s.auth.ValidateToken(ctx, token)
	if err != nil {
		slog.Error(fmt.Sprintf("validate token error: %s", err.Error()))
		return
	}

	if !validRole(role, api.UserRoleEmployee) && !validRole(role, api.UserRoleModerator) {
		slog.Error(fmt.Sprintf("bad role: %s", role))
		return
	}

	startDate := time.Time{}
	if params.StartDate != nil {
		startDate = *params.StartDate
	}

	endDate := time.Time{}
	if params.EndDate != nil {
		endDate = *params.EndDate
	}

	page := 1
	if params.Page != nil {
		page = *params.Page
	}

	limit := 10
	if params.Limit != nil {
		limit = *params.Limit
	}

	filter := dto.PvzInfoFilterDto{
		StartDate:      startDate,
		StartDateGiven: startDate != time.Time{},
		EndDate:        endDate,
		EndDateGiven:   endDate != time.Time{},
		Page:           page,
		Limit:          limit,
	}

	if !filter.IsValidParams() {
		slog.Error("bad filter params")
		return
	}

	serviceRes, err := s.pvz.GetPvzInfo(ctx, filter)
	if err != nil {
		slog.Debug(fmt.Sprintf("get pvz info error: %s", err.Error()))
		return
	}

	res := make([]api.GetPvzResponse, 0, len(serviceRes))
	for _, pvz := range serviceRes {
		pvzInfo := api.GetPvzResponse{}

		pvzInfo.Pvz = pvz.Pvz.ToApiModel()
		pvzInfo.Receptions = make([]struct {
			Reception api.Reception `json:"reception"`
			Products  []api.Product `json:"products"`
		}, 0, len(pvzInfo.Receptions))

		for _, reception := range pvz.Receptions {
			receptionInfo := struct {
				Reception api.Reception `json:"reception"`
				Products  []api.Product `json:"products"`
			}{}

			receptionInfo.Reception = reception.Reception.ToApiModel()
			receptionInfo.Products = make([]api.Product, 0, len(reception.Products))

			for _, product := range reception.Products {
				receptionInfo.Products = append(receptionInfo.Products, product.ToApiModel())
			}

			pvzInfo.Receptions = append(pvzInfo.Receptions, receptionInfo)
		}

		res = append(res, pvzInfo)
	}

	sendInfoResponse(w, res, http.StatusOK)
}
