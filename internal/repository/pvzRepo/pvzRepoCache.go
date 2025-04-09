package pvzRepo

import (
	"context"

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

func (r *PvzRepoCache) AddPvz(ctx context.Context, v model.Pvz) (model.Pvz, error) {
	r.arr = append(r.arr, v)

	return v, nil
}
