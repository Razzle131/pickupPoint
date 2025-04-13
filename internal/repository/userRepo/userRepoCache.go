package userRepo

import (
	"context"
	"errors"

	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/model"
)

type UserRepoCache struct {
	users []model.User
}

func NewCache() *UserRepoCache {
	return &UserRepoCache{
		users: make([]model.User, 0, consts.SliceMinCap),
	}
}

func (r *UserRepoCache) AddUser(ctx context.Context, user model.User) (model.User, error) {
	for _, u := range r.users {
		if u.Id == user.Id || u.Email == user.Email {
			return model.User{}, errors.New("not unique id")
		}
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
