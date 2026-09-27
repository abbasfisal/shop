package category

import (
	"context"
	"shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/interfaces/http/requests/admin"

	"github.com/gin-gonic/gin"
)

type CategoryServiceInterface interface {
	Index(ctx context.Context) (*responses.Categories, domain_err.CustomError)
	GetAllCategories(ctx context.Context) (*responses.Categories, domain_err.CustomError)
	GetAllParentCategory(ctx context.Context) (*responses.Categories, domain_err.CustomError)
	Show(ctx context.Context, categoryID int) (*responses.Category, domain_err.CustomError)
	CheckSlugUniqueness(ctx context.Context, slug string) bool
	Create(ctx context.Context, req *requests.CreateCategoryRequest) (*responses.Category, error)
	Edit(c *gin.Context, categoryID int, req *requests.UpdateCategoryRequest) domain_err.CustomError
}
