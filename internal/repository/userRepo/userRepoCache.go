package userRepo

import (
	"context"
	"errors"

	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
)

type UserRepoCache struct {
	users []model.User
}

func NewCache() *UserRepoCache {
	return &UserRepoCache{
		users: make([]model.User, 0, 16),
	}
}

// TODO: add ctx timeout
func (r *UserRepoCache) AddUser(ctx context.Context, dto dto.UserDto) (model.User, error) {
	user := model.User{
		Id:       uuid.New(),
		Email:    dto.Email,
		Password: dto.Password,
		Role:     dto.Role,
	}

	r.users = append(r.users, user)

	return user, nil
}

func (r *UserRepoCache) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return model.User{}, errors.New("not found")
}
