package model

import (
	"time"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/google/uuid"
)

type Reception struct {
	Id     uuid.UUID
	Date   time.Time
	PvzId  uuid.UUID
	Status api.ReceptionStatus
}
