package category

import (
	"context"
	"shop/internal/entities"
)

type CategoryRepositoryInterface interface {
	GetAll(ctx context.Context) ([]entities.Category, error)
	GetAllParent(ctx context.Context) ([]entities.Category, error)
	SelectBy(ctx context.Context, categoryID int) (entities.Category, error)
	FindBy(ctx context.Context, columnName string, value any) (entities.Category, error)
	Store(ctx context.Context, cat entities.Category) (entities.Category, error)
}
