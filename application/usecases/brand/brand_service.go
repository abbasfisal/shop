package brand

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
)

type BrandService struct {
	repo repositories.BrandRepositoryInterface
}

func NewBrandService(repo repositories.BrandRepositoryInterface) BrandServiceInterface {
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

func (bs *BrandService) Index(ctx context.Context) (*responses.Brands, domain_err.CustomError) {

	brands, err := bs.repo.GetAll(ctx)
	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToBrands(brands), domain_err.CustomError{}
}

func (bs *BrandService) Show(ctx context.Context, brandID int) (*responses.Brand, domain_err.CustomError) {

	fetchedBrand, err := bs.repo.SelectBy(ctx, brandID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_err.New(err.Error(), domain_err.RecordNotFound, 404)
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, domain_err.New(err.Error(), domain_err.MustBeUnique, 1062)
		}
		return nil, domain_err.New(err.Error(), domain_err.InternalServerError, 500)
	}

	return responses.ToBrand(fetchedBrand), domain_err.CustomError{}
}

func (bs *BrandService) Update(c *gin.Context, brandID int, req *requests.UpdateBrandRequest) (*responses.Brand, domain_err.CustomError) {

	updatedBrand, err := bs.repo.Update(c, brandID, req)
	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToBrand(updatedBrand), domain_err.CustomError{}
}
