package reception

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/Razzle131/pickupPoint/internal/repository/productRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/receptionRepo"
	"github.com/google/uuid"
)

type ReceptionService struct {
	products   productRepo.ProductRepo
	receptions receptionRepo.ReceptionRepo
}

func New(pr productRepo.ProductRepo, rp receptionRepo.ReceptionRepo) *ReceptionService {
	return &ReceptionService{
		products:   pr,
		receptions: rp,
	}
}

func (s *ReceptionService) CloseReception(pvzId uuid.UUID) error {
	reception, err := s.receptions.GetActiveReceptionByPvzId(context.TODO(), pvzId)
	if err != nil {
		return errors.New("open reception not found")
	}

	err = s.receptions.UpdateReceptionStatusById(context.Background(), reception.Id, string(api.Close))
	if err != nil {
		return errors.New("failed to update reception status")
	}

	return nil
}

// TODO: мб переделать на сортировку в бд и выдачу чисто последнего продукта, чтобы не гонять много данных
func (s *ReceptionService) RemoveReceptionLastProduct(pvzId uuid.UUID) error {
	reception, err := s.receptions.GetActiveReceptionByPvzId(context.TODO(), pvzId)
	if err != nil {
		return errors.New("open reception not found")
	}

	products, err := s.products.GetProductsByReceptionId(context.TODO(), reception.Id)
	if err != nil {
		return errors.New("GetProductsByReceptionId error")
	}

	if len(products) == 0 {
		return errors.New("nothing to delete")
	}

	// desc sort by date
	slices.SortFunc(products, func(a, b model.Product) int {
		diff := a.Date.Sub(b.Date)
		if diff < 0 {
			return 1
		}
		if diff > 0 {
			return -1
		}
		return 0
	})

	slog.Debug(fmt.Sprint(products))

	err = s.products.DeleteProductById(context.TODO(), products[0].Id)
	if err != nil {
		return errors.New("failed to delete product")
	}

	return nil
}

func (s *ReceptionService) AddReception(pvzId uuid.UUID) (model.Reception, error) {
	_, err := s.receptions.GetActiveReceptionByPvzId(context.TODO(), pvzId)
	if err == nil {
		return model.Reception{}, errors.New("close previous reception first")
	}

	reception, err := s.receptions.AddReception(context.TODO(), model.Reception{
		Id:     uuid.New(),
		Date:   time.Now(),
		PvzId:  pvzId,
		Status: api.InProgress,
	})
	if err != nil {
		return model.Reception{}, errors.New("failed to add new reception")
	}

	return reception, nil
}

func (s *ReceptionService) AddProduct(productType string, pvzId uuid.UUID) (model.Product, error) {
	reception, err := s.receptions.GetActiveReceptionByPvzId(context.TODO(), pvzId)
	if err != nil {
		return model.Product{}, errors.New("no receptions are in progress")
	}

	product, err := s.products.AddProduct(context.TODO(), model.Product{
		Id:          uuid.New(),
		Date:        time.Now(),
		Type:        api.ProductType(productType),
		ReceptionId: reception.Id,
	})

	if err != nil {
		return model.Product{}, err
	}

	return product, nil
}
