package pvzRepo

import (
	"context"

	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
)

type PvzRepo interface {
	AddPvz(ctx context.Context, pvz model.Pvz) (model.Pvz, error)
	GetPvzById(ctx context.Context, pvzId uuid.UUID) (model.Pvz, error)
	ListPvz(ctx context.Context) ([]model.Pvz, error)
}
