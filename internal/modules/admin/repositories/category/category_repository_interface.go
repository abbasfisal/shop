package category

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
	"shop/internal/modules/admin/requests"
)

type CategoryRepositoryInterface interface {
	GetAll(ctx context.Context) ([]*entities.Category, error)
	GetAllParent(ctx context.Context) ([]*entities.Category, error)
	SelectBy(ctx context.Context, categoryID int) (*entities.Category, error)
	FindBy(ctx context.Context, columnName string, value any) (*entities.Category, error)
	Store(ctx context.Context, cat *entities.Category) (*entities.Category, error)
	Update(c *gin.Context, categoryID int, req *requests.UpdateCategoryRequest) (*entities.Category, error)
}
