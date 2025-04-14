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

	// отдавать на откуп репозиторию не очень хотелось,
	// хотя по идее он должен следить за целостностью данных
	// напишите в фидбеке пж)
	_, err := s.pvzRepo.GetPvzById(ctx, pvz.Id)
	if err == nil {
		return dto.PvzDto{}, errors.New("pvz with such id already exists")
	}

	repoRes, err := s.pvzRepo.AddPvz(ctx, modelPvz)
	if err != nil {
		return dto.PvzDto{}, err
	}

	res := dto.PvzDto{}
	res.FromModel(repoRes)

	return res, nil
}

func (s *PvzService) GetPvzInfo(ctx context.Context, params dto.PvzInfoFilterDto) ([]dto.PvzInfoDto, error) {
	repoRes, err := s.pvzInfoRepo.ListPvzInfo(ctx, params)
	if err != nil {
		return []dto.PvzInfoDto{}, errors.New("failed to get pvzs info")
	}

	res := make([]dto.PvzInfoDto, 0, params.Limit)
	for pvz, receptions := range repoRes {
		pvzInfo := dto.PvzInfoDto{}

		pvzInfo.Pvz.FromModel(pvz)
		pvzInfo.Receptions = make([]dto.ReceptionInfoDto, 0, len(receptions))

		for reception, products := range receptions {
			receptionInfo := dto.ReceptionInfoDto{}

			receptionInfo.Reception.FromModel(reception)
			receptionInfo.Products = make([]dto.ProductDto, 0, len(products))

			for _, product := range products {
				var productDto dto.ProductDto
				productDto.FromModel(product)

				receptionInfo.Products = append(receptionInfo.Products, productDto)
			}

			pvzInfo.Receptions = append(pvzInfo.Receptions, receptionInfo)
		}

		res = append(res, pvzInfo)
	}

	return res, nil
}
