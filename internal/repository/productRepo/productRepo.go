package productRepo

import (
	"context"

	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
)

type ProductRepo interface {
	AddProduct(ctx context.Context, product model.Product) (model.Product, error)
	GetProductsByReceptionId(ctx context.Context, receptionId uuid.UUID) ([]model.Product, error)
	DeleteProductById(ctx context.Context, productId uuid.UUID) error
}
