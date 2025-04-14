package pvzInfoRepo

import (
	"testing"
	"time"

	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/Razzle131/pickupPoint/internal/repository/productRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/receptionRepo"
	"github.com/google/uuid"
)

func TestListPvzInfo(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	receptionRepo := receptionRepo.NewCache(pvzRepo)
	prodRepo := productRepo.NewCache(receptionRepo)

	repo := NewCache(pvzRepo, prodRepo, receptionRepo)

	params := dto.PvzInfoFilterDto{
		StartDateGiven: false,
		EndDateGiven:   false,
		Page:           1,
		Limit:          1,
	}

	res, err := repo.ListPvzInfo(t.Context(), params)
	if err != nil {
		t.Errorf("list pvz info error: %s", err.Error())
	}

	if len(res) != 0 {
		t.Errorf("bad result")
	}

	pvz := model.Pvz{
		Id:      uuid.New(),
		City:    "aaa",
		RegDate: time.Now(),
	}

	reception := model.Reception{
		Id:     uuid.New(),
		Date:   time.Now(),
		PvzId:  pvz.Id,
		Status: consts.ReceptionStatusActive,
	}

	product := model.Product{
		Id:          uuid.New(),
		Date:        time.Now(),
		Type:        consts.ProductTypeClothes,
		ReceptionId: reception.Id,
	}

	repo.pvzRepo.AddPvz(t.Context(), pvz)
	repo.receptionRepo.AddReception(t.Context(), reception)
	repo.productRepo.AddProduct(t.Context(), product)

	res, err = repo.ListPvzInfo(t.Context(), params)
	if err != nil {
		t.Errorf("list pvz info error: %s", err.Error())
	}

	if len(res) != 1 {
		t.Errorf("bad result")
	}

	if _, found := res[pvz]; !found {
		t.Errorf("bad result")
	}

	products, found := res[pvz][reception]
	if !found {
		t.Errorf("bad result")
	}

	if len(products) != 1 {
		t.Errorf("bad products len: %v", len(products))
	}

	params.Limit = 1
	res, err = repo.ListPvzInfo(t.Context(), params)
	if err != nil {
		t.Errorf("list pvz info error with limit 1: %s", err.Error())
	}

	if len(res) != 1 {
		t.Errorf("bad result")
	}

	if _, found := res[pvz]; !found {
		t.Errorf("bad result")
	}

	products, found = res[pvz][reception]
	if !found {
		t.Errorf("bad result")
	}

	if len(products) != 1 {
		t.Errorf("bad products len: %v", len(products))
	}

	params.Page = 20
	res, err = repo.ListPvzInfo(t.Context(), params)
	if err != nil {
		t.Errorf("list pvz info error with limit 1 and page 20: %s", err.Error())
	}

	if len(res) != 0 {
		t.Errorf("bad result")
	}

	params.Page = 1
	params.Limit = 1
	params.EndDate = time.Now().Add(time.Hour)
	params.EndDateGiven = true
	res, err = repo.ListPvzInfo(t.Context(), params)
	if err != nil {
		t.Errorf("list pvz info error with limit 1 and page 20: %s", err.Error())
	}

	if len(res) != 1 {
		t.Errorf("bad result")
	}
}
