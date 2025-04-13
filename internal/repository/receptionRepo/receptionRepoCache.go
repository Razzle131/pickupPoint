package receptionRepo

import (
	"context"
	"errors"

	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzRepo"
	"github.com/google/uuid"
)

type ReceptionRepoCache struct {
	receptions []model.Reception
	pvzRepo    pvzRepo.PvzRepo
}

func NewCache(pvzRepo pvzRepo.PvzRepo) *ReceptionRepoCache {
	return &ReceptionRepoCache{
		receptions: make([]model.Reception, 0, consts.SliceMinCap),
		pvzRepo:    pvzRepo,
	}
}

func (r *ReceptionRepoCache) AddReception(ctx context.Context, reception model.Reception) (model.Reception, error) {
	_, err := r.pvzRepo.GetPvzById(context.Background(), reception.PvzId)
	if err != nil {
		return model.Reception{}, errors.New("cant find such pvz")
	}

	for _, r := range r.receptions {
		if r.Id == reception.Id {
			return model.Reception{}, errors.New("not unique id")
		}
	}

	r.receptions = append(r.receptions, reception)

	return reception, nil
}

func (r *ReceptionRepoCache) GetActiveReceptionByPvzId(ctx context.Context, pvzId uuid.UUID) (model.Reception, error) {
	for _, reception := range r.receptions {
		if reception.PvzId == pvzId && reception.Status == consts.ReceptionStatusActive {
			return reception, nil
		}
	}

	return model.Reception{}, errors.New("no active receptions")
}

func (r *ReceptionRepoCache) GetReceptionById(ctx context.Context, id uuid.UUID) (model.Reception, error) {
	for _, reception := range r.receptions {
		if reception.Id == id {
			return reception, nil
		}
	}

	return model.Reception{}, errors.New("not found")
}

func (r *ReceptionRepoCache) UpdateReceptionStatusById(ctx context.Context, receptionId uuid.UUID, newStatus string) (model.Reception, error) {
	for i := 0; i < len(r.receptions); i++ {
		if r.receptions[i].Id == receptionId {
			r.receptions[i].Status = newStatus
			return r.receptions[i], nil
		}
	}

	return model.Reception{}, errors.New("not found")
}

func (r *ReceptionRepoCache) GetReceptionsByPvzId(ctx context.Context, pvzId uuid.UUID) ([]model.Reception, error) {
	res := make([]model.Reception, 0, consts.SliceMinCap)

	for _, reception := range r.receptions {
		if reception.PvzId == pvzId {
			res = append(res, reception)
		}
	}

	return res, nil
}
