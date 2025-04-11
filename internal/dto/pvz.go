package dto

import (
	"time"

	"github.com/google/uuid"
)

type PvzDto struct {
	Id      uuid.UUID
	City    string
	RegDate time.Time
}
