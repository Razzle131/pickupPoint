package receptionRepo

import (
	"testing"
	"time"

	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzRepo"
	"github.com/google/uuid"
)

func TestAddReception(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	repo := NewCache(pvzRepo)

	reception := model.Reception{
		Id:     uuid.New(),
		Date:   time.Now(),
		PvzId:  uuid.New(),
		Status: consts.ReceptionStatusActive,
	}

	_, err := repo.AddReception(t.Context(), reception)
	if err == nil {
		t.Errorf("error must occur when trying to add reception to non existing pvz")
	}

	repo.pvzRepo.AddPvz(t.Context(), model.Pvz{
		Id:      reception.PvzId,
		City:    "abc",
		RegDate: time.Now(),
	})

	res, err := repo.AddReception(t.Context(), reception)
	if err != nil {
		t.Errorf("add reception error: %s", err.Error())
	}

	if res != reception {
		t.Errorf("bad result data")
	}

	if repo.receptions[0] != reception {
		t.Errorf("bad written data")
	}

	_, err = repo.AddReception(t.Context(), reception)
	if err == nil {
		t.Errorf("error must occur when adding non unique reception id")
	}
}

func TestGetActiveReceptionByPvzId(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	repo := NewCache(pvzRepo)

	reception := model.Reception{
		Id:     uuid.New(),
		Date:   time.Now(),
		PvzId:  uuid.New(),
		Status: consts.ReceptionStatusActive,
	}

	_, err := repo.GetActiveReceptionByPvzId(t.Context(), reception.PvzId)
	if err == nil {
		t.Errorf("error must occur when no active reception found")
	}

	repo.pvzRepo.AddPvz(t.Context(), model.Pvz{
		Id:      reception.PvzId,
		City:    "abc",
		RegDate: time.Now(),
	})
	repo.AddReception(t.Context(), reception)

	res, err := repo.GetActiveReceptionByPvzId(t.Context(), reception.PvzId)
	if err != nil {
		t.Errorf("get active reception error: %s", err.Error())
	}

	if res != reception {
		t.Errorf("bad result data")
	}
}

func TestGetReceptionById(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	repo := NewCache(pvzRepo)

	reception := model.Reception{
		Id:     uuid.New(),
		Date:   time.Now(),
		PvzId:  uuid.New(),
		Status: consts.ReceptionStatusActive,
	}

	_, err := repo.GetReceptionById(t.Context(), reception.Id)
	if err == nil {
		t.Errorf("error must occur when no reception found")
	}

	repo.pvzRepo.AddPvz(t.Context(), model.Pvz{Id: reception.PvzId})
	repo.AddReception(t.Context(), reception)

	res, err := repo.GetReceptionById(t.Context(), reception.Id)
	if err != nil {
		t.Errorf("get reception error: %s", err.Error())
	}

	if res != reception {
		t.Errorf("bad result data")
	}
}

func TestUpdateReceptionStatusById(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	repo := NewCache(pvzRepo)

	reception := model.Reception{
		Id:     uuid.New(),
		Date:   time.Now(),
		PvzId:  uuid.New(),
		Status: consts.ReceptionStatusActive,
	}

	_, err := repo.UpdateReceptionStatusById(t.Context(), reception.Id, consts.ReceptionStatusClosed)
	if err == nil {
		t.Errorf("error must occur when no reception found")
	}

	repo.pvzRepo.AddPvz(t.Context(), model.Pvz{
		Id:      reception.PvzId,
		City:    "abc",
		RegDate: time.Now(),
	})
	repo.AddReception(t.Context(), reception)

	res, err := repo.UpdateReceptionStatusById(t.Context(), reception.Id, consts.ReceptionStatusClosed)
	if err != nil {
		t.Errorf("update reception error: %s", err.Error())
	}

	if res.Id != reception.Id ||
		res.Date != reception.Date ||
		res.PvzId != reception.PvzId ||
		res.Status != consts.ReceptionStatusClosed {
		t.Errorf("bad returned data")
	}

	if repo.receptions[0].Id != reception.Id ||
		repo.receptions[0].Date != reception.Date ||
		repo.receptions[0].PvzId != reception.PvzId ||
		repo.receptions[0].Status != consts.ReceptionStatusClosed {
		t.Errorf("bad written data")
	}
}

func TestGetReceptionsByPvzId(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	repo := NewCache(pvzRepo)

	reception := model.Reception{
		Id:     uuid.New(),
		Date:   time.Now(),
		PvzId:  uuid.New(),
		Status: consts.ReceptionStatusActive,
	}

	res, _ := repo.GetReceptionsByPvzId(t.Context(), reception.PvzId)
	if len(res) != 0 {
		t.Errorf("bad result data")
	}

	repo.pvzRepo.AddPvz(t.Context(), model.Pvz{
		Id:      reception.PvzId,
		City:    "abc",
		RegDate: time.Now(),
	})
	repo.AddReception(t.Context(), reception)

	res, _ = repo.GetReceptionsByPvzId(t.Context(), reception.PvzId)
	if len(res) != 1 {
		t.Errorf("bad result data")
	}

	if res[0] != reception {
		t.Errorf("bad result data")
	}
}
