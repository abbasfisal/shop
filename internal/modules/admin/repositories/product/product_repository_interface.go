package product

import (
	"context"
	"shop/internal/entities"
)

type ProductRepositoryInterface interface {
	GetAll(ctx context.Context) ([]entities.Product, error)
	FindBy(ctx context.Context, columnName string, value any) (entities.Product, error)
	Store(ctx context.Context, product entities.Product) (entities.Product, error)
}
