package home

import (
	"context"
	"shop/internal/entities"
)

type HomeRepositoryInterface interface {
	GetRandomProducts(ctx context.Context, limit int) ([]entities.Product, error)
	GetLatestProducts(ctx context.Context, limit int) ([]entities.Product, error)
	GetCategories(ctx context.Context, limit int) ([]entities.Category, error)
	GetProduct(ctx context.Context, productSlug, sku string) (entities.Product, error)
	GetProductsBy(ctx context.Context, columnName string, value any) ([]entities.Product, error)
	GetCategoryBy(ctx context.Context, columnName string, value any) (entities.Category, error)
}
