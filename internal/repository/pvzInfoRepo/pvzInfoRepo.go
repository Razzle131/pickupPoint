package pvzInfoRepo

import (
	"context"

	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/model"
)

// For avoiding N+1 queries problem when we want to list pvz`s info
type PvzInfoRepo interface {
	ListPvzInfo(ctx context.Context, params dto.PvzInfoFilterDto) (map[model.Pvz]map[model.Reception][]model.Product, error)
}
