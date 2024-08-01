package brand

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type BrandServiceInterface interface {
	CheckSlugUniqueness(ctx context.Context, slug string) bool
	Create(ctx context.Context, req requests.CreateBrandRequest) (responses.Brand, error)
	Index(ctx context.Context) (responses.Brands, custom_error.CustomError)
	Show(ctx context.Context, brandID int) (responses.Brand, custom_error.CustomError)
	Update(c *gin.Context, brandID int, req requests.UpdateBrandRequest) (responses.Brand, custom_error.CustomError)
}
