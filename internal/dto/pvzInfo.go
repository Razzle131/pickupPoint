package dto

import (
	"github.com/Razzle131/pickupPoint/internal/model"
)

type ReceptionInfoDto struct {
	Reception model.Reception
	Products  []model.Product
}

type PvzInfoDto struct {
	Pvz        model.Pvz
	Receptions []ReceptionInfoDto
}
