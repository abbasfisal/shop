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

func NewBrandService(repo brand.BrandRepositoryInterface) BrandServiceInterface {
	return &BrandService{repo: repo}
}

func (bs *BrandService) CheckSlugUniqueness(ctx context.Context, slug string) bool {
	existingBrand, _ := bs.repo.FindBy(ctx, "slug", slug)
	if existingBrand.ID > 0 {
		return true
	}
	return false
}

func (bs *BrandService) Create(ctx context.Context, req *requests.CreateBrandRequest) (*responses.Brand, error) {

	var brandToCreate = entities.Brand{
		Title: req.Title,
		Slug:  req.Slug,
		Image: req.Image,
	}

	newBrand, err := bs.repo.Store(ctx, &brandToCreate)
	if err != nil {
		return nil, err
	}

	return responses.ToBrand(newBrand), nil
}

func (bs *BrandService) Index(ctx context.Context) (*responses.Brands, custom_error.CustomError) {

	brands, err := bs.repo.GetAll(ctx)
	if err != nil {
		return nil, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToBrands(brands), custom_error.CustomError{}
}

func (bs *BrandService) Show(ctx context.Context, brandID int) (*responses.Brand, custom_error.CustomError) {

	fetchedBrand, err := bs.repo.SelectBy(ctx, brandID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, custom_error.New(err.Error(), custom_error.MustBeUnique, 1062)
		}
		return nil, custom_error.New(err.Error(), custom_error.InternalServerError, 500)
	}

	return responses.ToBrand(fetchedBrand), custom_error.CustomError{}
}

func (bs *BrandService) Update(c *gin.Context, brandID int, req *requests.UpdateBrandRequest) (*responses.Brand, custom_error.CustomError) {

	updatedBrand, err := bs.repo.Update(c, brandID, req)
	if err != nil {
		return nil, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToBrand(updatedBrand), custom_error.CustomError{}
}
