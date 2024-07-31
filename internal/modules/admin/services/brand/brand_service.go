package brand

import (
	"context"
	"shop/internal/entities"
	"shop/internal/modules/admin/repositories/brand"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
)

type BrandService struct {
	repo brand.BrandRepositoryInterface
}

func NewBrandService(repo brand.BrandRepositoryInterface) *BrandService {
	return &BrandService{repo: repo}
}

func (bs BrandService) CheckSlugUniqueness(ctx context.Context, slug string) bool {
	brand, _ := bs.repo.FindBy(ctx, "slug", slug)
	if brand.ID > 0 {
		return true
	}
	return false
}

func (bs BrandService) Create(ctx context.Context, req requests.CreateBrandRequest) (responses.Brand, error) {

	var response responses.Brand

	var brand = entities.Brand{
		Title: req.Title,
		Slug:  req.Slug,
		Image: req.Image,
	}

	newBrand, err := bs.repo.Store(ctx, brand)
	if err != nil {
		return response, err
	}

	return responses.ToBrand(newBrand), nil
}
