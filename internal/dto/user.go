package dto

import (
	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

type UserDto struct {
	Id    uuid.UUID
	Email string
	Role  string
}

func (d UserDto) IsValidRole() bool {
	return d.Role == consts.UserRoleEmployee || d.Role == consts.UserRoleModerator
}

func (d *UserDto) FromModel(u model.User) {
	d.Email = u.Email
	d.Id = u.Id
	d.Role = u.Role
}

func (d UserDto) ToModel() model.User {
	return model.User{
		Id:       d.Id,
		Email:    d.Email,
		Password: "",
		Role:     d.Role,
	}
}

func (d *UserDto) FromApiModel(p api.User) {
	id := uuid.New()
	if p.Id != nil {
		id = *p.Id
	}

	d.Email = string(p.Email)
	d.Id = id
	d.Role = string(p.Role)
}

func (d UserDto) ToApiModel() api.User {
	return api.User{
		Id:    &d.Id,
		Email: types.Email(d.Email),
		Role:  api.UserRole(d.Role),
	}
}
