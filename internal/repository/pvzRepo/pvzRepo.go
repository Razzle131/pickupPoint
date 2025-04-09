package pvzRepo

import (
	"context"

	"github.com/Razzle131/pickupPoint/internal/model"
)

type PvzRepo interface {
	AddPvz(ctx context.Context, v model.Pvz) (model.Pvz, error)
}
