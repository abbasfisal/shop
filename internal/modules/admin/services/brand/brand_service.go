package brand

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/internal/entities"
	"shop/internal/modules/admin/repositories/brand"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
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

	var brand = entities.Brand{
		Title: req.Title,
		Slug:  req.Slug,
		Image: req.Image,
	}

	newBrand, err := bs.repo.Store(ctx, brand)
	if err != nil {
		return responses.Brand{}, err
	}

	return responses.ToBrand(newBrand), nil
}

func (bs BrandService) Index(ctx context.Context) (responses.Brands, custom_error.CustomError) {

	brands, err := bs.repo.GetAll(ctx)
	if err != nil {
		return responses.Brands{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToBrands(brands), custom_error.CustomError{}
}

func (bs BrandService) Show(ctx context.Context, brandID int) (responses.Brand, custom_error.CustomError) {
	var response responses.Brand

	brand, err := bs.repo.SelectBy(ctx, brandID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return response, custom_error.New(err.Error(), custom_error.MustBeUnique, 1062)
		}
		return response, custom_error.New(err.Error(), custom_error.InternalServerError, 500)
	}

	return responses.ToBrand(brand), custom_error.CustomError{}
}

func (bs BrandService) Update(c *gin.Context, brandID int, req requests.UpdateBrandRequest) (responses.Brand, custom_error.CustomError) {

	brand, err := bs.repo.Update(c, brandID, req)
	if err != nil {
		return responses.Brand{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToBrand(brand), custom_error.CustomError{}
}
