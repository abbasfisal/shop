package brand

import (
	"context"
	"shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/interfaces/http/requests/admin"

	"github.com/gin-gonic/gin"
)

type BrandServiceInterface interface {
	CheckSlugUniqueness(ctx context.Context, slug string) bool
	Create(ctx context.Context, req *requests.CreateBrandRequest) (*responses.Brand, error)
	Index(ctx context.Context) (*responses.Brands, domain_err.CustomError)
	Show(ctx context.Context, brandID int) (*responses.Brand, domain_err.CustomError)
	Update(c *gin.Context, brandID int, req *requests.UpdateBrandRequest) (*responses.Brand, domain_err.CustomError)
}
