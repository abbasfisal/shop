package category

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"shop/internal/entities"
	"shop/internal/modules/admin/repositories/category"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type CategoryService struct {
	repo category.CategoryRepositoryInterface
}

func NewCategoryService() CategoryService {
	return CategoryService{
		repo: category.NewCategoryRepository(),
	}
}
func (cs CategoryService) Index(ctx context.Context) (responses.Categories, custom_error.CustomError) {

	var response responses.Categories
	categories, err := cs.repo.GetAll(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		return response, custom_error.New(err.Error(), custom_error.InternalServerError, 500)
	}

	return responses.ToCategories(categories), custom_error.CustomError{}
}

func (cs CategoryService) Show(ctx context.Context, categoryID int) (responses.Category, custom_error.CustomError) {
	var response responses.Category

	cat, err := cs.repo.SelectBy(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		return response, custom_error.New(err.Error(), custom_error.InternalServerError, 500)
	}

	return responses.ToCategory(cat), custom_error.CustomError{}
}

func (cs CategoryService) CheckSlugUniqueness(ctx context.Context, slug string) bool {
	cat, _ := cs.repo.FindBy(ctx, "slug", slug)
	if cat.ID > 0 {
		return true
	}
	return false
}

func (cs CategoryService) Create(ctx context.Context, req requests.CreateCategoryRequest) (responses.Category, error) {

	var response responses.Category

	var cat = entities.Category{
		Title: req.Title,
		Slug:  req.Slug,
		Image: req.Image,
		Status: func() bool {
			if req.Status == "on" {
				return true
			}
			return false
		}(),
	}
	newCategory, err := cs.repo.Store(ctx, cat)
	if err != nil {
		return response, err
	}

	return responses.ToCategory(newCategory), nil
}

func (cs CategoryService) GetAllCategories(ctx context.Context) (responses.Categories, custom_error.CustomError) {
	var response responses.Categories
	categories, err := cs.repo.GetAll(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		return response, custom_error.New(err.Error(), custom_error.InternalServerError, 500)
	}

	return responses.ToCategories(categories), custom_error.CustomError{}
}
