package category

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
)

type CategoryService struct {
	repo repositories.CategoryRepositoryInterface
}

func NewCategoryService(categoryRepo repositories.CategoryRepositoryInterface) CategoryServiceInterface {
	return &CategoryService{
		repo: categoryRepo,
	}
}

func (cs *CategoryService) Index(ctx context.Context) (*responses.Categories, domain_err.CustomError) {

	categories, err := cs.repo.GetAll(ctx)
	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToCategories(categories), domain_err.CustomError{}
}

func (cs *CategoryService) Show(ctx context.Context, categoryID int) (*responses.Category, domain_err.CustomError) {

	cat, err := cs.repo.SelectBy(ctx, categoryID)
	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToCategory(cat), domain_err.CustomError{}
}

func (cs *CategoryService) CheckSlugUniqueness(ctx context.Context, slug string) bool {
	cat, _ := cs.repo.FindBy(ctx, "slug", slug)
	if cat.ID > 0 {
		return true
	}
	return false
}

func (cs *CategoryService) Create(ctx context.Context, req *requests.CreateCategoryRequest) (*responses.Category, error) {

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
	if req.ParentID != 0 {
		cat.ParentID = &req.ParentID
	}

	if req.Priority != nil {
		if *req.Priority != 0 {
			cat.Priority = req.Priority
		}
	}
	newCategory, err := cs.repo.Store(ctx, &cat)
	if err != nil {
		return nil, err
	}
	return responses.ToCategory(newCategory), nil
}

func (cs *CategoryService) GetAllCategories(ctx context.Context) (*responses.Categories, domain_err.CustomError) {

	categories, err := cs.repo.GetAll(ctx)
	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToCategories(categories), domain_err.CustomError{}
}

func (cs *CategoryService) GetAllParentCategory(ctx context.Context) (*responses.Categories, domain_err.CustomError) {

	categories, err := cs.repo.GetAllParent(ctx)
	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToCategories(categories), domain_err.CustomError{}
}

func (cs *CategoryService) Edit(c *gin.Context, categoryID int, req *requests.UpdateCategoryRequest) domain_err.CustomError {
	_, err := cs.repo.Update(c, categoryID, req)

	if err != nil {
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return domain_err.CustomError{}
}
