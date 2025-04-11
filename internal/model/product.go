package model

import (
	"time"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID
	Date        time.Time
	Type        api.ProductType
	ReceptionId uuid.UUID
}
