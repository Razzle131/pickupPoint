package model

import (
	"github.com/Razzle131/pickupPoint/api"
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID
	Email    string
	Password string
	Role     api.UserRole
}
