package dto

import (
	"time"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
)

type PvzDto struct {
	Id      uuid.UUID
	City    string
	RegDate time.Time
}

func (d PvzDto) IsValidCity() bool {
	return d.City == consts.PvzCityKazan || d.City == consts.PvzCityMoscow || d.City == consts.PvzCityStPetersburg
}

func (d *PvzDto) FromModel(p model.Pvz) {
	d.City = p.City
	d.Id = p.Id
	d.RegDate = p.RegDate
}

func (d PvzDto) ToModel() model.Pvz {
	return model.Pvz{
		Id:      d.Id,
		City:    d.City,
		RegDate: d.RegDate,
	}
}

func (d *PvzDto) FromApiModel(p api.PVZ) {
	date := time.Now()
	if p.RegistrationDate != nil {
		date = *p.RegistrationDate
	}

	id := uuid.New()
	if p.Id != nil {
		id = *p.Id
	}

	d.Id = id
	d.City = string(p.City)
	d.RegDate = date
}

func (d PvzDto) ToApiModel() api.PVZ {
	return api.PVZ{
		City:             api.PVZCity(d.City),
		Id:               &d.Id,
		RegistrationDate: &d.RegDate,
	}
}
