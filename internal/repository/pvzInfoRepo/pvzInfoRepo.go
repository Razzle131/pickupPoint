package pvzInfoRepo

import (
	"context"

	"github.com/Razzle131/pickupPoint/internal/dto"
)

// For avoiding N+1 queries problem when we want to list all pvz`s info
type PvzInfoRepo interface {
	ListPvzInfo(ctx context.Context, params dto.PvzInfoFilterDto) ([]dto.PvzInfoDto, error)
}
