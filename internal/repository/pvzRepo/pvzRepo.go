package pvzRepo

import (
	"context"

	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/model"
)

// TODO: подумать о замене апи модели на дто (вопрос в том что в чем смысл этого)
type PvzRepo interface {
	AddPvz(ctx context.Context, dto dto.PvzDto) (model.Pvz, error)
	ListPvz(ctx context.Context, params dto.PvzInfoFilter) ([]model.Pvz, error)
}
