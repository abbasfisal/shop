package home

import (
	"context"
	"shop/internal/modules/admin/responses"
	"shop/internal/modules/public/repositories/home"
	"shop/internal/pkg/custom_error"
)

type HomeService struct {
	repo home.HomeRepositoryInterface
}

func NewHomeService() HomeService {
	return HomeService{
		repo: home.NewHomeRepository(),
	}
}

func (h HomeService) GetProducts(ctx context.Context, limit int) (responses.Products, custom_error.CustomError) {

	products, err := h.repo.GetLatestProducts(ctx, limit)
	if err != nil {
		return responses.Products{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToProducts(products), custom_error.CustomError{}
}

func (h HomeService) GetCategories(ctx context.Context, limit int) (responses.Categories, custom_error.CustomError) {

	categories, err := h.repo.GetCategories(ctx, limit)
	if err != nil {
		return responses.Categories{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToCategories(categories), custom_error.CustomError{}
}

func (h HomeService) ShowProductDetail(ctx context.Context, productSlug, sku string) (responses.Product, custom_error.CustomError) {

	product, err := h.repo.GetProduct(ctx, productSlug, sku)
	if err != nil {
		return responses.Product{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToProduct(product), custom_error.CustomError{}
}

func (h HomeService) ShowProductsByCategorySlug(ctx context.Context, value any) (responses.Products, custom_error.CustomError) {

	//get
	category, cErr := h.ShowCategory(ctx, "slug", value)
	if cErr.Code == 404 {
		return responses.Products{}, custom_error.New(cErr.Error(), custom_error.RecordNotFound, 404)
	}
	if cErr.Code == 500 {
		return responses.Products{}, custom_error.New(cErr.Error(), custom_error.InternalServerError, 500)
	}

	products, err := h.repo.GetProductsBy(ctx, "category_id", category.ID)
	if err != nil {
		return responses.Products{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToProducts(products), custom_error.CustomError{}
}

func (h HomeService) ShowCategory(ctx context.Context, columnName string, value any) (responses.Category, custom_error.CustomError) {

	category, err := h.repo.GetCategoryBy(ctx, columnName, value)
	if err != nil {
		return responses.Category{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToCategory(category), custom_error.CustomError{}
}
