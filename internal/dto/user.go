package dto

import "github.com/Razzle131/pickupPoint/api"

type UserDto struct {
	Email    string
	Password string
	Role     api.UserRole
}
