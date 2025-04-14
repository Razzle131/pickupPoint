package pvzRepo

import (
	"testing"
	"time"

	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
)

func TestAddPvz(t *testing.T) {
	repo := NewCache()

	pvz := model.Pvz{
		Id:      uuid.New(),
		City:    "abc",
		RegDate: time.Now(),
	}

	res, err := repo.AddPvz(t.Context(), pvz)
	if err != nil {
		t.Errorf("add pvz error: %s", err.Error())
	}

	if res.Id != pvz.Id || res.City != pvz.City || res.RegDate != pvz.RegDate {
		t.Errorf("bad result data")
	}

	if repo.arr[0].Id != pvz.Id || repo.arr[0].City != pvz.City || repo.arr[0].RegDate != pvz.RegDate {
		t.Errorf("bad written data")
	}

	_, err = repo.AddPvz(t.Context(), pvz)
	if err == nil {
		t.Errorf("error must occur when adding non unique pvz")
	}
}

func TestGetPvzById(t *testing.T) {
	repo := NewCache()

	pvz := model.Pvz{
		Id:      uuid.New(),
		City:    "abc",
		RegDate: time.Now(),
	}

	_, err := repo.GetPvzById(t.Context(), pvz.Id)
	if err == nil {
		t.Errorf("error must occur when no pvz found")
	}

	repo.AddPvz(t.Context(), pvz)

	res, err := repo.GetPvzById(t.Context(), pvz.Id)
	if err != nil {
		t.Errorf("get pvz error: %s", err.Error())
	}

	if res.Id != pvz.Id || res.City != pvz.City || res.RegDate != pvz.RegDate {
		t.Errorf("bad result data")
	}
}

func TestListPvz(t *testing.T) {
	repo := NewCache()

	pvz := model.Pvz{
		Id:      uuid.New(),
		City:    "abc",
		RegDate: time.Now(),
	}

	res, err := repo.ListPvz(t.Context())
	if err != nil {
		t.Errorf("list pvz error: %s", err.Error())
	}

	if len(res) != 0 {
		t.Errorf("bad result data")
	}

	repo.AddPvz(t.Context(), pvz)

	res, err = repo.ListPvz(t.Context())
	if err != nil {
		t.Errorf("list pvz error: %s", err.Error())
	}

	if len(res) != 1 {
		t.Errorf("bad result data")
	}

	if res[0].Id != pvz.Id || res[0].City != pvz.City || res[0].RegDate != pvz.RegDate {
		t.Errorf("given data doesnt much input")
	}
}
