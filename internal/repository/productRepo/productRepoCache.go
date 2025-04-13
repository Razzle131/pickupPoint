package productRepo

import (
	"context"
	"errors"
	"slices"

	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/Razzle131/pickupPoint/internal/repository/receptionRepo"
	"github.com/google/uuid"
)

type ProductRepoCache struct {
	products   []model.Product
	receptions receptionRepo.ReceptionRepo
}

func NewCache(receptions receptionRepo.ReceptionRepo) *ProductRepoCache {
	return &ProductRepoCache{
		products:   make([]model.Product, 0, consts.SliceMinCap),
		receptions: receptions,
	}
}

func (r *ProductRepoCache) AddProduct(ctx context.Context, product model.Product) (model.Product, error) {
	// smt like foreign key validation
	if _, err := r.receptions.GetReceptionById(context.Background(), product.ReceptionId); err != nil {
		return model.Product{}, errors.New("reception not found")
	}

	for _, p := range r.products {
		if p.Id == product.Id {
			return model.Product{}, errors.New("not unique id")
		}
	}

	r.products = append(r.products, product)

	return product, nil
}

func (r *ProductRepoCache) DeleteProductById(ctx context.Context, productId uuid.UUID) error {
	for i := 0; i < len(r.products); i++ {
		if r.products[i].Id == productId {
			r.products = slices.Delete(r.products, i, i+1)
			return nil
		}
	}

	return errors.New("not found")
}

func (r *ProductRepoCache) GetProductsByReceptionId(ctx context.Context, receptionId uuid.UUID) ([]model.Product, error) {
	res := make([]model.Product, 0, consts.SliceMinCap)

	for _, product := range r.products {
		if product.ReceptionId == receptionId {
			res = append(res, product)
		}
	}

	return res, nil
}

func (r *ProductRepoCache) GetReceptionLastProduct(ctx context.Context, receptionId uuid.UUID) (model.Product, error) {
	res := model.Product{}
	foundOne := false

	for _, product := range r.products {
		if product.ReceptionId == receptionId && res.Date.Sub(product.Date) < 0 {
			foundOne = true
			res = product
		}
	}

	if !foundOne {
		return model.Product{}, errors.New("not found")
	}

	return res, nil
}
