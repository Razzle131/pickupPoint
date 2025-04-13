package pvzInfoRepo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/repository/productRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/receptionRepo"
)

type PvzInfoRepoCache struct {
	pvzRepo       pvzRepo.PvzRepo
	productRepo   productRepo.ProductRepo
	receptionRepo receptionRepo.ReceptionRepo
}

func NewCache(pvzRepo pvzRepo.PvzRepo, productRepo productRepo.ProductRepo, receptionRepo receptionRepo.ReceptionRepo) *PvzInfoRepoCache {
	return &PvzInfoRepoCache{
		pvzRepo:       pvzRepo,
		productRepo:   productRepo,
		receptionRepo: receptionRepo,
	}
}

// TODO: переделать пагинацию и вывод данных (не должно быть пвз без приемок)
func (r *PvzInfoRepoCache) ListPvzInfo(ctx context.Context, params dto.PvzInfoFilterDto) ([]dto.PvzInfoDto, error) {
	Pvzs, err := r.pvzRepo.ListPvz(context.TODO(), params)
	if err != nil {
		return []dto.PvzInfoDto{}, errors.New("failed to list pvzs")
	}

	res := make([]dto.PvzInfoDto, 0, len(Pvzs))
	for _, pvz := range Pvzs {
		receptions, err := r.receptionRepo.GetReceptionsByPvzId(context.TODO(), pvz.Id)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to get receptions for pvz id %s", pvz.Id))
			continue
		}

		receptionsInfo := make([]dto.ReceptionInfoDto, 0, len(receptions))
		for _, reception := range receptions {
			if (params.StartDateGiven && reception.Date.Sub(params.StartDate) < 0) ||
				(params.EndDateGiven && reception.Date.Sub(params.EndDate) > 0) {
				continue
			}

			products, err := r.productRepo.GetProductsByReceptionId(context.TODO(), reception.Id)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			receptionsInfo = append(receptionsInfo, dto.ReceptionInfoDto{
				Reception: reception,
				Products:  products,
			})
		}

		res = append(res, dto.PvzInfoDto{
			Pvz:        pvz,
			Receptions: receptionsInfo,
		})
	}

	return res, nil
}
