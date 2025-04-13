package reception

import (
	"context"
	"errors"
	"time"

	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/Razzle131/pickupPoint/internal/repository/productRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/receptionRepo"
	"github.com/google/uuid"
)

type ReceptionService struct {
	products   productRepo.ProductRepo
	receptions receptionRepo.ReceptionRepo
	pvzs       pvzRepo.PvzRepo
}

func New(pr productRepo.ProductRepo, rp receptionRepo.ReceptionRepo, pvzs pvzRepo.PvzRepo) *ReceptionService {
	return &ReceptionService{
		products:   pr,
		receptions: rp,
		pvzs:       pvzs,
	}
}

func (s *ReceptionService) CloseReception(ctx context.Context, pvzId uuid.UUID) (dto.ReceptionDto, error) {
	reception, err := s.receptions.GetActiveReceptionByPvzId(ctx, pvzId)
	if err != nil {
		return dto.ReceptionDto{}, errors.New("open reception not found")
	}

	repoRes, err := s.receptions.UpdateReceptionStatusById(ctx, reception.Id, consts.ReceptionStatusClosed)
	if err != nil {
		return dto.ReceptionDto{}, errors.New("failed to update reception status")
	}

	res := dto.ReceptionDto{}
	res.FromModel(repoRes)

	return res, nil
}

func (s *ReceptionService) RemoveReceptionLastProduct(ctx context.Context, pvzId uuid.UUID) (dto.ProductDto, error) {
	reception, err := s.receptions.GetActiveReceptionByPvzId(ctx, pvzId)
	if err != nil {
		return dto.ProductDto{}, errors.New("open reception not found")
	}

	product, err := s.products.GetReceptionLastProduct(ctx, reception.Id)
	if err != nil {
		return dto.ProductDto{}, errors.New("GetProductsByReceptionId error")
	}

	err = s.products.DeleteProductById(ctx, product.Id)
	if err != nil {
		return dto.ProductDto{}, errors.New("failed to delete product")
	}

	res := dto.ProductDto{}
	res.FromModel(product)

	return res, nil
}

func (s *ReceptionService) AddReception(ctx context.Context, pvzId uuid.UUID) (dto.ReceptionDto, error) {
	_, err := s.pvzs.GetPvzById(ctx, pvzId)
	if err != nil {
		return dto.ReceptionDto{}, errors.New("no such pvz")
	}

	_, err = s.receptions.GetActiveReceptionByPvzId(ctx, pvzId)
	if err == nil {
		return dto.ReceptionDto{}, errors.New("close previous reception first")
	}

	reception, err := s.receptions.AddReception(ctx, model.Reception{
		Id:     uuid.New(),
		Date:   time.Now(),
		PvzId:  pvzId,
		Status: consts.ReceptionStatusActive,
	})
	if err != nil {
		return dto.ReceptionDto{}, errors.New("failed to add new reception")
	}

	res := dto.ReceptionDto{}
	res.FromModel(reception)

	return res, nil
}

func (s *ReceptionService) AddProduct(ctx context.Context, productType string, pvzId uuid.UUID) (dto.ProductDto, error) {
	reception, err := s.receptions.GetActiveReceptionByPvzId(ctx, pvzId)
	if err != nil {
		return dto.ProductDto{}, errors.New("no receptions are in progress")
	}

	product, err := s.products.AddProduct(ctx, model.Product{
		Id:          uuid.New(),
		Date:        time.Now(),
		Type:        productType,
		ReceptionId: reception.Id,
	})
	if err != nil {
		return dto.ProductDto{}, err
	}

	res := dto.ProductDto{}
	res.FromModel(product)

	return res, nil
}
