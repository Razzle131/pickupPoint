package pvzRepo

import (
	"context"

	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/model"
)

type PvzRepoCache struct {
	arr []model.Pvz
}

func NewCache() *PvzRepoCache {
	return &PvzRepoCache{
		arr: make([]model.Pvz, 0, 16),
	}
}

func (r *PvzRepoCache) AddPvz(ctx context.Context, dto dto.PvzDto) (model.Pvz, error) {
	pvz := model.Pvz{
		Id:      dto.Id,
		City:    dto.City,
		RegDate: dto.RegDate,
	}

	r.arr = append(r.arr, pvz)

	return pvz, nil
}

func (r *PvzRepoCache) ListPvz(ctx context.Context, params dto.PvzInfoFilter) ([]model.Pvz, error) {
	start := (params.Page - 1) * params.Limit
	end := start + params.Limit

	res := make([]model.Pvz, 0, params.Limit)
	for i := min(start, len(r.arr)); i < min(end, len(r.arr)); i++ {
		res = append(res, r.arr[i])
	}

	return res, nil
}
