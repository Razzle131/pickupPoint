package receptionRepo

import (
	"context"
	"errors"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
)

type ReceptionRepoCache struct {
	receptions []model.Reception
}

func NewCache() *ReceptionRepoCache {
	return &ReceptionRepoCache{
		receptions: make([]model.Reception, 0, 16),
	}
}

func (r *ReceptionRepoCache) AddReception(ctx context.Context, reception model.Reception) (model.Reception, error) {
	r.receptions = append(r.receptions, reception)

	return reception, nil
}

func (r *ReceptionRepoCache) GetActiveReceptionByPvzId(ctx context.Context, pvzId uuid.UUID) (model.Reception, error) {
	for _, reception := range r.receptions {
		if reception.Id == pvzId && reception.Status == api.InProgress {
			return reception, nil
		}
	}

	return model.Reception{}, errors.New("no active receptions")
}

func (r *ReceptionRepoCache) UpdateReceptionStatusById(ctx context.Context, receptionId uuid.UUID, newStatus string) error {
	for i := 0; i < len(r.receptions); i++ {
		if r.receptions[i].Id == receptionId {
			r.receptions[i].Status = api.ReceptionStatus(newStatus)
			return nil
		}
	}

	return errors.New("not found")
}

func (r *ReceptionRepoCache) GetReceptionsByPvzId(ctx context.Context, pvzId uuid.UUID) ([]model.Reception, error) {
	res := make([]model.Reception, 0, 16)

	for _, reception := range r.receptions {
		if reception.PvzId == pvzId {
			res = append(res, reception)
		}
	}

	return res, nil
}
