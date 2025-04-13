package dto

import (
	"time"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
)

type ProductDto struct {
	Id          uuid.UUID
	Date        time.Time
	Type        string
	ReceptionId uuid.UUID
}

func (d ProductDto) IsValidType() bool {
	return d.Type == consts.ProductTypeDevices || d.Type == consts.ProductTypeShoes || d.Type == consts.ProductTypeClothes
}

func (d *ProductDto) FromModel(p model.Product) {
	d.Date = p.Date
	d.Id = p.Id
	d.ReceptionId = p.ReceptionId
	d.Type = p.Type
}

func (d ProductDto) ToModel() model.Product {
	return model.Product{
		Id:          d.Id,
		Date:        d.Date,
		Type:        d.Type,
		ReceptionId: d.ReceptionId,
	}
}

func (d *ProductDto) FromApiModel(p api.Product) {
	date := time.Now()
	if p.DateTime != nil {
		date = *p.DateTime
	}

	id := uuid.New()
	if p.Id != nil {
		id = *p.Id
	}

	d.Date = date
	d.Id = id
	d.ReceptionId = p.ReceptionId
	d.Type = string(p.Type)
}

func (d ProductDto) ToApiModel() api.Product {
	return api.Product{
		DateTime:    &d.Date,
		Id:          &d.Id,
		ReceptionId: d.ReceptionId,
		Type:        api.ProductType(d.Type),
	}
}
