package home

import (
	"context"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type HomeServiceInterface interface {
	GetProducts(ctx context.Context, limit int) (responses.Products, custom_error.CustomError)
	GetCategories(ctx context.Context, limit int) (responses.Categories, custom_error.CustomError)
	ShowCategory(ctx context.Context, columnName string, value any) (responses.Category, custom_error.CustomError)
	ShowProductDetail(ctx context.Context, productSlug, sku string) (responses.Product, custom_error.CustomError)
	ShowProductsByCategorySlug(ctx context.Context, value any) (responses.Products, custom_error.CustomError)
}
