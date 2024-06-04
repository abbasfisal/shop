package product

import (
	"context"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type ProductServiceInterface interface {
	Index(ctx context.Context) (responses.Products, custom_error.CustomError)
	Show(ctx context.Context, columnName string, value any) (responses.Product, custom_error.CustomError)
	Create(ctx context.Context, req requests.CreateProductRequest) (responses.Product, custom_error.CustomError)
	CheckSkuIsUnique(ctx context.Context, sku string) (bool, custom_error.CustomError)
}
