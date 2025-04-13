package pvz

import (
	"testing"
	"time"

	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/repository/productRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzInfoRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/receptionRepo"
	"github.com/google/uuid"
)

func TestCreatePvz(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	service := New(pvzRepo, nil)

	id := uuid.New()
	city := "Москва"
	regDate := time.Now()

	res, err := service.CreatePvz(t.Context(), dto.PvzDto{
		Id:      id,
		City:    city,
		RegDate: regDate,
	})

	if err != nil {
		t.Errorf("create pvz error: %s", err.Error())
	}

	if res.Id != id || res.City != city || res.RegDate != regDate {
		t.Errorf("wrong result fields")
	}

	_, err = service.CreatePvz(t.Context(), dto.PvzDto{
		Id:      id,
		City:    city,
		RegDate: regDate,
	})

	if err == nil {
		t.Errorf("error must occur when adding same pvz")
	}
}

func TestGetPvzInfo(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	receptionRepo := receptionRepo.NewCache(pvzRepo)
	prodRepo := productRepo.NewCache(receptionRepo)
	pvzInfoRepo := pvzInfoRepo.NewCache(pvzRepo, prodRepo, receptionRepo)

	service := New(pvzRepo, pvzInfoRepo)

	_, err := service.GetPvzInfo(t.Context(), dto.PvzInfoFilterDto{Page: 1, Limit: 1})
	if err != nil {
		t.Errorf("empty get pvz info error: %s", err.Error())
	}

	id := uuid.New()
	city := "Москва"
	regDate := time.Now()

	service.CreatePvz(t.Context(), dto.PvzDto{
		Id:      id,
		City:    city,
		RegDate: regDate,
	})

	res, err := service.GetPvzInfo(t.Context(), dto.PvzInfoFilterDto{
		StartDateGiven: false,
		EndDateGiven:   false,
		Page:           1,
		Limit:          10,
	})
	if err != nil {
		t.Errorf("get pvz info error: %s", err.Error())
	}

	if len(res) != 1 {
		t.Errorf("bad get pvz info res len: %v", len(res))
	}

	if res[0].Pvz.Id != id || len(res[0].Receptions) > 0 {
		t.Errorf("bad get pvz info response")
	}

	// no res by filters check
	res, err = service.GetPvzInfo(t.Context(), dto.PvzInfoFilterDto{
		StartDateGiven: false,
		EndDateGiven:   true,
		Page:           20,
		Limit:          10,
	})
	if err != nil {
		t.Errorf("get pvz info error: %s", err.Error())
	}

	if len(res) != 0 {
		t.Errorf("bad get pvz info res len: %v", len(res))
	}
}
