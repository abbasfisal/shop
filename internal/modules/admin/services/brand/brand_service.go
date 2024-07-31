package brand

import (
	"shop/internal/modules/admin/repositories/brand"
)

type BrandService struct {
	repo brand.BrandRepositoryInterface
}

func NewBrandService(repo brand.BrandRepositoryInterface) *BrandService {
	return &BrandService{repo: repo}
}
