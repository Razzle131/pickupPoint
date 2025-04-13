package dto

import (
	"time"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
)

type ReceptionDto struct {
	Id       uuid.UUID
	DateTime time.Time
	PvzId    uuid.UUID
	Status   string
}

func (d ReceptionDto) IsValidStatus() bool {
	return d.Status == consts.ReceptionStatusActive || d.Status == consts.ReceptionStatusClosed
}

func (d *ReceptionDto) FromModel(r model.Reception) {
	d.DateTime = r.Date
	d.Id = r.Id
	d.PvzId = r.PvzId
	d.Status = r.Status
}

func (d ReceptionDto) ToModel() model.Reception {
	return model.Reception{
		Id:     d.Id,
		Date:   d.DateTime,
		PvzId:  d.PvzId,
		Status: d.Status,
	}
}

func (d *ReceptionDto) FromApiModel(p api.Reception) {
	id := uuid.New()
	if p.Id != nil {
		id = *p.Id
	}

	d.DateTime = p.DateTime
	d.Id = id
	d.PvzId = p.PvzId
	d.Status = string(p.Status)
}

func (d ReceptionDto) ToApiModel() api.Reception {
	return api.Reception{
		Id:       &d.Id,
		DateTime: d.DateTime,
		PvzId:    d.PvzId,
		Status:   api.ReceptionStatus(d.Status),
	}
}
