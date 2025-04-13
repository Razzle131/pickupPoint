package userRepo

import (
	"testing"

	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
)

func TestAddUser(t *testing.T) {
	repo := NewCache()

	id := uuid.New()
	user := model.User{
		Id:       id,
		Email:    "abc@abc.abc",
		Password: "abc",
		Role:     "moderator",
	}

	res, err := repo.AddUser(t.Context(), user)
	if err != nil {
		t.Errorf("add user error: %s", err.Error())
	}

	if res != user {
		t.Errorf("bad response")
	}

	if repo.users[0] != user {
		t.Errorf("bad written result")
	}

	user = model.User{
		Id:       id,
		Email:    "abccc@abc.abc",
		Password: "abc",
		Role:     "moderator",
	}

	_, err = repo.AddUser(t.Context(), user)
	if err == nil {
		t.Errorf("error must occur when adding non unique user id")
	}

	user = model.User{
		Id:       uuid.New(),
		Email:    "abc@abc.abc",
		Password: "abc",
		Role:     "moderator",
	}

	_, err = repo.AddUser(t.Context(), user)
	if err == nil {
		t.Errorf("error must occur when adding non unique email")
	}
}

func TestGetUserByEmail(t *testing.T) {
	repo := NewCache()

	id := uuid.New()
	user := model.User{
		Id:       id,
		Email:    "abc@abc.abc",
		Password: "abc",
		Role:     "moderator",
	}

	_, err := repo.GetUserByEmail(t.Context(), user.Email)
	if err == nil {
		t.Errorf("error must occur when trying to get unregistered user")
	}

	repo.AddUser(t.Context(), user)

	res, err := repo.GetUserByEmail(t.Context(), user.Email)
	if err != nil {
		t.Errorf("get user by email error: %s", err.Error())
	}

	if res.Email != user.Email || res.Id != user.Id || res.Password != user.Password || res.Role != user.Role {
		t.Errorf("bad result data")
	}
}
