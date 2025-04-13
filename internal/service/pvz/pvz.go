package pvz

import (
	"context"
	"errors"

	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzInfoRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzRepo"
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

func (s *PvzService) CreatePvz(ctx context.Context, pvz dto.PvzDto) (dto.PvzDto, error) {
	modelPvz := pvz.ToModel()

	repoRes, err := s.pvzRepo.AddPvz(ctx, modelPvz)
	if err != nil {
		return dto.PvzDto{}, err
	}

	res := dto.PvzDto{}
	res.FromModel(repoRes)

	return res, nil
}

func (s *PvzService) GetPvzInfo(ctx context.Context, params dto.PvzInfoFilterDto) ([]dto.PvzInfoDto, error) {
	res, err := s.pvzInfoRepo.ListPvzInfo(ctx, params)
	if err != nil {
		return []dto.PvzInfoDto{}, errors.New("failed to get pvzs info")
	}

	return res, nil
}
