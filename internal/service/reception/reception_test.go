package reception

import (
	"testing"

	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/Razzle131/pickupPoint/internal/repository/productRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/receptionRepo"
	"github.com/google/uuid"
)

func TestCloseReception(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	receptionRepo := receptionRepo.NewCache(pvzRepo)
	prodRepo := productRepo.NewCache(receptionRepo)

	service := New(prodRepo, receptionRepo, pvzRepo)

	id := uuid.New()
	_, err := service.CloseReception(t.Context(), id)
	if err == nil {
		t.Errorf("error must occur when no receptions active")
	}

	service.pvzs.AddPvz(t.Context(), model.Pvz{Id: id})
	reception, _ := service.AddReception(t.Context(), id)

	res, err := service.CloseReception(t.Context(), id)
	if err != nil {
		t.Errorf("close reception error: %s", err.Error())
	}

	if reception.Id != res.Id {
		t.Errorf("closed wrong reception")
	}

	_, err = service.CloseReception(t.Context(), id)
	if err == nil {
		t.Errorf("error must occur when no receptions are active")
	}
}

func TestRemoveReceptionLastProduct(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	receptionRepo := receptionRepo.NewCache(pvzRepo)
	prodRepo := productRepo.NewCache(receptionRepo)

	service := New(prodRepo, receptionRepo, pvzRepo)

	pvzId := uuid.New()
	_, err := service.RemoveReceptionLastProduct(t.Context(), pvzId)
	if err == nil {
		t.Errorf("error must occur when no receptions are active")
	}

	service.pvzs.AddPvz(t.Context(), model.Pvz{Id: pvzId})
	service.AddReception(t.Context(), pvzId)

	_, err = service.RemoveReceptionLastProduct(t.Context(), pvzId)
	if err == nil {
		t.Errorf("error must occur when no products in reception")
	}

	service.AddProduct(t.Context(), consts.ProductTypeShoes, pvzId)
	service.AddProduct(t.Context(), consts.ProductTypeDevices, pvzId)

	prod, err := service.RemoveReceptionLastProduct(t.Context(), pvzId)
	if err != nil {
		t.Errorf("error on removing product")
	}

	if prod.Type != consts.ProductTypeDevices {
		t.Errorf("removed wrong product")
	}
}

func TestAddReception(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	receptionRepo := receptionRepo.NewCache(pvzRepo)
	prodRepo := productRepo.NewCache(receptionRepo)

	service := New(prodRepo, receptionRepo, pvzRepo)

	pvzId := uuid.New()
	_, err := service.AddReception(t.Context(), pvzId)
	if err == nil {
		t.Errorf("error must occur when no pvz presented")
	}

	service.pvzs.AddPvz(t.Context(), model.Pvz{Id: pvzId})

	reception, err := service.AddReception(t.Context(), pvzId)
	if err != nil {
		t.Errorf("add reception error: %s", err.Error())
	}

	if reception.PvzId != pvzId || reception.Status != consts.ReceptionStatusActive {
		t.Errorf("bad reception created")
	}

	_, err = service.AddReception(t.Context(), pvzId)
	if err == nil {
		t.Errorf("error must occur when adding reception without closing previous")
	}

	service.CloseReception(t.Context(), pvzId)

	_, err = service.AddReception(t.Context(), pvzId)
	if err != nil {
		t.Errorf("add reception error: %s", err.Error())
	}
}

func TestAddProduct(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	receptionRepo := receptionRepo.NewCache(pvzRepo)
	prodRepo := productRepo.NewCache(receptionRepo)

	service := New(prodRepo, receptionRepo, pvzRepo)

	pvzId := uuid.New()

	_, err := service.AddProduct(t.Context(), consts.ProductTypeClothes, pvzId)
	if err == nil {
		t.Errorf("error must occur when no pvz and reception")
	}

	service.pvzs.AddPvz(t.Context(), model.Pvz{Id: pvzId})

	_, err = service.AddProduct(t.Context(), consts.ProductTypeClothes, pvzId)
	if err == nil {
		t.Errorf("error must occur when no reception are open")
	}

	reception, _ := service.AddReception(t.Context(), pvzId)

	product, err := service.AddProduct(t.Context(), consts.ProductTypeClothes, pvzId)
	if err != nil {
		t.Errorf("add product error: %s", err.Error())
	}

	if product.Type != consts.ProductTypeClothes || product.ReceptionId != reception.Id {
		t.Errorf("bad data returned")
	}

	service.CloseReception(t.Context(), pvzId)

	_, err = service.AddProduct(t.Context(), consts.ProductTypeClothes, pvzId)
	if err == nil {
		t.Errorf("error must occur when previous reception closed")
	}
}
