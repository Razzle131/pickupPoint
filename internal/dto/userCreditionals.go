package dto

import (
	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/model"
)

type UserCreditionalsDto struct {
	Email    string
	Password string
	Role     string
}

func (d UserCreditionalsDto) IsValidRole() bool {
	return d.Role == consts.UserRoleEmployee || d.Role == consts.UserRoleModerator
}

func (d *UserCreditionalsDto) FromModel(u model.User) {
	d.Email = u.Email
	d.Password = u.Password
	d.Role = u.Role
}

func (d UserCreditionalsDto) ToModel() model.User {
	return model.User{
		Email:    d.Email,
		Password: d.Password,
		Role:     d.Role,
	}
}
