package model

import (
	"time"

	"github.com/google/uuid"
)

type Pvz struct {
	Id      uuid.UUID
	City    string
	RegDate time.Time
}
