package receptionRepo

import (
	"context"

	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
)

type ReceptionRepo interface {
	AddReception(ctx context.Context, reception model.Reception) (model.Reception, error)
	GetActiveReceptionByPvzId(ctx context.Context, pvzId uuid.UUID) (model.Reception, error)
	GetReceptionsByPvzId(ctx context.Context, pvzId uuid.UUID) ([]model.Reception, error)
	UpdateReceptionStatusById(ctx context.Context, receptionId uuid.UUID, newStatus string) error
}
