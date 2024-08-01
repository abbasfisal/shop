package brand

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
	"shop/internal/modules/admin/requests"
)

type BrandRepositoryInterface interface {
	GetAll(ctx context.Context) ([]entities.Brand, error)
	FindBy(ctx context.Context, columnName string, value any) (entities.Brand, error)
	Store(ctx context.Context, brand entities.Brand) (entities.Brand, error)
	SelectBy(ctx context.Context, brandID int) (entities.Brand, error)
	Update(c *gin.Context, brandID int, req requests.UpdateBrandRequest) (entities.Brand, error)
}
