package userRepo

import (
	"context"

	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/model"
)

type UserRepo interface {
	AddUser(ctx context.Context, dto dto.UserDto) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}
