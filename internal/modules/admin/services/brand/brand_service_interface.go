package brand

import (
	"context"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
)

type BrandServiceInterface interface {
	CheckSlugUniqueness(ctx context.Context, slug string) bool
	Create(ctx context.Context, req requests.CreateBrandRequest) (responses.Brand, error)
}
