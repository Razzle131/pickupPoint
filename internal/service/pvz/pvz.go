package pvz

import (
	"context"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzRepo"
	"github.com/oapi-codegen/runtime/types"
)

type PvzService struct {
	pvzRepo pvzRepo.PvzRepo
}

func New(repo pvzRepo.PvzRepo) *PvzService {
	return &PvzService{
		pvzRepo: repo,
	}
}

// TODO: мб перенести проверку роли сюда, т.к. роль это по факту не чек входа
func (s *PvzService) CreatePvz(req api.PVZ) (api.PVZ, error) {
	model := model.Pvz{
		Id:      *req.Id,
		City:    string(req.City),
		RegDate: *req.RegistrationDate,
	}

	repoRes, err := s.pvzRepo.AddPvz(context.TODO(), model)
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
