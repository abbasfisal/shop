package category

import (
	"context"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type CategoryServiceInterface interface {
	Index(ctx context.Context) (responses.Categories, custom_error.CustomError)
	GetAllCategories(ctx context.Context) (responses.Categories, custom_error.CustomError)
	Show(ctx context.Context, categoryID int) (responses.Category, custom_error.CustomError)
	CheckSlugUniqueness(ctx context.Context, slug string) bool
	Create(ctx context.Context, req requests.CreateCategoryRequest) (responses.Category, error)
}
