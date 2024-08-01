package brand

import (
	"context"
	"shop/internal/entities"
)

type BrandRepositoryInterface interface {
	GetAll(ctx context.Context) ([]entities.Brand, error)
	FindBy(ctx context.Context, columnName string, value any) (entities.Brand, error)
	Store(ctx context.Context, brand entities.Brand) (entities.Brand, error)
	SelectBy(ctx context.Context, brandID int) (entities.Brand, error)
}
