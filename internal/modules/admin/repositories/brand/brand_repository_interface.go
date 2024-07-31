package brand

import (
	"context"
	"shop/internal/entities"
)

type BrandRepositoryInterface interface {
	FindBy(ctx context.Context, columnName string, value any) (entities.Brand, error)
	Store(ctx context.Context, brand entities.Brand) (entities.Brand, error)
}
