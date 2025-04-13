package productRepo

import (
	"testing"
	"time"

	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/receptionRepo"
	"github.com/google/uuid"
)

func TestAddProduct(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	receptionRepo := receptionRepo.NewCache(pvzRepo)
	repo := NewCache(receptionRepo)

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

	_, err := repo.AddProduct(t.Context(), product)
	if err == nil {
		t.Errorf("error must occur when adding product to nonexisting reception")
	}

	pvzRepo.AddPvz(t.Context(), pvz)
	receptionRepo.AddReception(t.Context(), reception)

	res, err := repo.AddProduct(t.Context(), product)
	if err != nil {
		t.Errorf("add product error: %s", err.Error())
	}

	if res != product {
		t.Errorf("unexpected result")
	}

	_, err = repo.AddProduct(t.Context(), product)
	if err == nil {
		t.Errorf("error must occur when adding non unique product")
	}

	// no test for adding to closed reception because it is business rule
	// (adding to non existing reception breaks data consistency, so it is tested)
}

func TestDeleteProductById(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	receptionRepo := receptionRepo.NewCache(pvzRepo)
	repo := NewCache(receptionRepo)

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

	err := repo.DeleteProductById(t.Context(), product.Id)
	if err == nil {
		t.Errorf("error must occur on deleting non existing product")
	}

	pvzRepo.AddPvz(t.Context(), pvz)
	receptionRepo.AddReception(t.Context(), reception)
	repo.AddProduct(t.Context(), product)

	err = repo.DeleteProductById(t.Context(), product.Id)
	if err != nil {
		t.Errorf("delete product error: %s", err.Error())
	}
}

func TestGetProductsByReceptionId(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	receptionRepo := receptionRepo.NewCache(pvzRepo)
	repo := NewCache(receptionRepo)

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

	res, err := repo.GetProductsByReceptionId(t.Context(), product.ReceptionId)
	if err != nil {
		t.Errorf("get products error: %s", err.Error())
	}

	if len(res) > 0 {
		t.Errorf("unexpected result")
	}

	pvzRepo.AddPvz(t.Context(), pvz)
	receptionRepo.AddReception(t.Context(), reception)
	repo.AddProduct(t.Context(), product)

	res, err = repo.GetProductsByReceptionId(t.Context(), product.ReceptionId)
	if err != nil {
		t.Errorf("get products error: %s", err.Error())
	}

	if len(res) != 1 || res[0] != product {
		t.Errorf("unexpected result")
	}
}

func TestGetReceptionLastProduct(t *testing.T) {
	pvzRepo := pvzRepo.NewCache()
	receptionRepo := receptionRepo.NewCache(pvzRepo)
	repo := NewCache(receptionRepo)

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

	_, err := repo.GetReceptionLastProduct(t.Context(), product.ReceptionId)
	if err == nil {
		t.Errorf("error must occur on getting product of non existing reception")
	}

	pvzRepo.AddPvz(t.Context(), pvz)
	receptionRepo.AddReception(t.Context(), reception)

	_, err = repo.GetReceptionLastProduct(t.Context(), product.ReceptionId)
	if err == nil {
		t.Errorf("error must occur on getting product of empty reception")
	}

	repo.AddProduct(t.Context(), product)

	res, err := repo.GetReceptionLastProduct(t.Context(), product.ReceptionId)
	if err != nil {
		t.Errorf("getting last product error: %s", err.Error())
	}

	if res != product {
		t.Errorf("unexpected result")
	}
}
