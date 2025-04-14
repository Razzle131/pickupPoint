package pvz

import (
	"testing"
	"time"

	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/model"
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

	pvz := dto.PvzDto{
		Id:      uuid.New(),
		City:    consts.PvzCityMoscow,
		RegDate: time.Now(),
	}

	service.CreatePvz(t.Context(), pvz)

	params := dto.PvzInfoFilterDto{
		StartDateGiven: false,
		EndDateGiven:   false,
		Page:           1,
		Limit:          10,
	}

	res, err := service.GetPvzInfo(t.Context(), params)
	if err != nil {
		t.Errorf("get pvz info error: %s", err.Error())
	}

	if len(res) != 0 {
		t.Errorf("bad get pvz info res len for pvz without receptions: %v", len(res))
	}

	receptionRepo.AddReception(t.Context(), model.Reception{
		Id:    uuid.New(),
		Date:  time.Now(),
		PvzId: pvz.Id,
	})

	res, err = service.GetPvzInfo(t.Context(), params)
	if err != nil {
		t.Errorf("get pvz info error: %s", err.Error())
	}

	if len(res) != 1 {
		t.Errorf("bad get pvz info res len for pvz with empty reception: %v", len(res))
	}

	// no res by filters check
	res, err = service.GetPvzInfo(t.Context(), dto.PvzInfoFilterDto{
		StartDateGiven: false,
		EndDate:        time.Time{},
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
