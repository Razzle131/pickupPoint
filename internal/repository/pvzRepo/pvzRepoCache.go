package pvzRepo

import (
	"context"
	"errors"

	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
)

type PvzRepoCache struct {
	arr []model.Pvz
}

func NewCache() *PvzRepoCache {
	return &PvzRepoCache{
		arr: make([]model.Pvz, 0, consts.SliceMinCap),
	}
}

func (r *PvzRepoCache) AddPvz(ctx context.Context, pvz model.Pvz) (model.Pvz, error) {
	for _, p := range r.arr {
		if p.Id == pvz.Id {
			return model.Pvz{}, errors.New("not unique id")
		}
	}

	r.arr = append(r.arr, pvz)

	return pvz, nil
}

func (r *PvzRepoCache) GetPvzById(ctx context.Context, pvzId uuid.UUID) (model.Pvz, error) {
	for _, pvz := range r.arr {
		if pvz.Id == pvzId {
			return pvz, nil
		}
	}

	return model.Pvz{}, errors.New("not found")
}

func (r *PvzRepoCache) ListPvz(ctx context.Context) ([]model.Pvz, error) {
	res := make([]model.Pvz, len(r.arr))
	copy(res, r.arr)

	return res, nil
}
