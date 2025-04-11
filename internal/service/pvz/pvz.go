package pvz

import (
	"context"
	"errors"
	"time"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzInfoRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzRepo"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

type PvzService struct {
	pvzRepo     pvzRepo.PvzRepo
	pvzInfoRepo pvzInfoRepo.PvzInfoRepo
}

func New(pvzRepo pvzRepo.PvzRepo, pvzInfoRepo pvzInfoRepo.PvzInfoRepo) *PvzService {
	return &PvzService{
		pvzRepo:     pvzRepo,
		pvzInfoRepo: pvzInfoRepo,
	}
}

// TODO: мб перенести проверку роли сюда, т.к. роль это по факту не чек входа
func (s *PvzService) CreatePvz(req api.PVZ) (api.PVZ, error) {
	reqId := uuid.New()
	reqDate := time.Now()

	if req.Id != nil {
		reqId = *req.Id
	}
	if req.RegistrationDate != nil {
		reqDate = *req.RegistrationDate
	}

	pvz := dto.PvzDto{
		Id:      reqId,
		City:    string(req.City),
		RegDate: reqDate,
	}

	repoRes, err := s.pvzRepo.AddPvz(context.TODO(), pvz)
	if err != nil {
		return api.PVZ{}, err
	}

	id := types.UUID(repoRes.Id)
	res := api.PVZ{
		City:             api.PVZCity(repoRes.City),
		Id:               &id,
		RegistrationDate: &repoRes.RegDate,
	}

	return res, nil
}

func (s *PvzService) GetPvzInfo(params dto.PvzInfoFilter) ([]dto.PvzInfoDto, error) {
	res, err := s.pvzInfoRepo.ListPvzInfo(context.TODO(), params)
	if err != nil {
		return []dto.PvzInfoDto{}, errors.New("failed to get pvzs info")
	}

	return res, nil
}
