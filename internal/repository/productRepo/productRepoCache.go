package productRepo

import (
	"context"
	"errors"
	"slices"

	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
)

type ProductRepoCache struct {
	products []model.Product
}

func NewCache() *ProductRepoCache {
	return &ProductRepoCache{
		products: make([]model.Product, 0, 16),
	}
}

func (r *ProductRepoCache) AddProduct(ctx context.Context, product model.Product) (model.Product, error) {
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
	res := make([]model.Product, 0, 16)

	for _, product := range r.products {
		if product.ReceptionId == receptionId {
			res = append(res, product)
		}
	}

	return res, nil
}
