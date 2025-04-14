package pvzInfoRepo

import (
	"context"
	"errors"

	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/model"
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

func (r *PvzInfoRepoCache) ListPvzInfo(ctx context.Context, params dto.PvzInfoFilterDto) (map[model.Pvz]map[model.Reception][]model.Product, error) {
	pvzs, err := r.pvzRepo.ListPvz(ctx)
	if err != nil {
		return map[model.Pvz]map[model.Reception][]model.Product{}, errors.New("failed to list pvzs")
	}

	res := make(map[model.Pvz](map[model.Reception][]model.Product), params.Limit)

	offset := (params.Page - 1) * params.Limit
	rows := 0

	for _, pvz := range pvzs {
		receptions, err := r.receptionRepo.GetReceptionsByPvzId(ctx, pvz.Id)
		if err != nil || len(receptions) == 0 {
			continue
		}

		if offset <= 0 {
			res[pvz] = make(map[model.Reception][]model.Product, params.Limit)
		}

		for _, reception := range receptions {
			if (params.StartDateGiven && reception.Date.Sub(params.StartDate) <= 0) ||
				(params.EndDateGiven && reception.Date.Sub(params.EndDate) >= 0) {
				continue
			}

			products, err := r.productRepo.GetProductsByReceptionId(ctx, reception.Id)
			if err != nil {
				continue
			}

			if len(products) == 0 {
				offset--
				if offset <= 0 {
					rows++
				}
			}

			if offset <= 0 {
				res[pvz][reception] = make([]model.Product, 0, len(products))
			}

			for _, product := range products {
				if offset > 0 {
					offset--
					continue
				}

				res[pvz][reception] = append(res[pvz][reception], product)
				rows++
			}
		}
	}

	return res, nil
}
