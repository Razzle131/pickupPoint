package authorization

import (
	"testing"

	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/repository/userRepo"
)

func TestDummyLogin(t *testing.T) {
	service := New(nil)

	role := consts.UserRoleModerator

	token, err := service.DummyLogin(t.Context(), role)
	if err != nil {
		t.Errorf("dummy login error: %s", err.Error())
	}

	parsedRole, err := service.ValidateToken(t.Context(), token)
	if err != nil {
		t.Errorf("dummy login error: %s", err.Error())
	}

	if parsedRole != role {
		t.Errorf("not expected role: %s", parsedRole)
	}
}

func TestLogin(t *testing.T) {
	repo := userRepo.NewCache()
	service := New(repo)

	email := "abc@abc.abc"
	password := "aboba"

	dto := dto.UserCreditionalsDto{
		Email:    email,
		Password: password,
	}

	_, err := service.Login(t.Context(), dto)

	if err == nil {
		t.Errorf("login error must occur")
	}

	service.Register(t.Context(), dto)

	_, err = service.Login(t.Context(), dto)
	if err != nil {
		t.Errorf("login error occurred: %s", err.Error())
	}
}

func TestRegister(t *testing.T) {
	repo := userRepo.NewCache()
	service := New(repo)

	email := "abc@abc.abc"
	password := "aboba"
	role := consts.UserRoleModerator

	user, err := service.Register(t.Context(), dto.UserCreditionalsDto{
		Email:    email,
		Password: password,
		Role:     role,
	})

	if err != nil {
		t.Errorf("register error: %s", err.Error())
	}

	if user.Email != email || user.Role != role {
		t.Errorf("not expected result parameters: %s %s", user.Email, user.Role)
	}

	_, err = service.Register(t.Context(), dto.UserCreditionalsDto{
		Email:    email,
		Password: password,
		Role:     role,
	})

	if err == nil {
		t.Errorf("register error must occur")
	}
}

func TestValidateToken(t *testing.T) {
	repo := userRepo.NewCache()
	service := New(repo)

	token, _ := service.DummyLogin(t.Context(), consts.UserRoleModerator)

	role, err := service.ValidateToken(t.Context(), token)
	if err != nil {
		t.Errorf("validate token error: %s", err.Error())
	}

	if role != consts.UserRoleModerator {
		t.Errorf("parsed bad role: %s", role)
	}

	_, err = service.ValidateToken(t.Context(), "abc")
	if err == nil {
		t.Errorf("error shoud occur with bad key")
	}

	_, err = service.ValidateToken(t.Context(), "bearer abc")
	if err == nil {
		t.Errorf("error shoud occur with bad key")
	}
}
