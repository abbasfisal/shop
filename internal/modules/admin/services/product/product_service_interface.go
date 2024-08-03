package product

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type ProductServiceInterface interface {
	Index(ctx context.Context) (responses.Products, custom_error.CustomError)
	Show(ctx context.Context, columnName string, value any) (responses.Product, custom_error.CustomError)
	Create(ctx context.Context, req requests.CreateProductRequest) (responses.Product, custom_error.CustomError)
	CheckSkuIsUnique(ctx context.Context, sku string) (bool, custom_error.CustomError)
	FetchByProductID(c *gin.Context, productID int) (responses.Product, custom_error.CustomError)
	FetchRootAttributes(c *gin.Context, productID int) (responses.Attributes, custom_error.CustomError)
	AddAttributeValues(c *gin.Context, productID int, attributes []string) custom_error.CustomError
	FetchProductAttributes(c *gin.Context, productID int) (responses.Product, custom_error.CustomError)
	CreateInventory(c *gin.Context, productID int, req requests.CreateProductInventoryRequest) custom_error.CustomError
	FetchImage(c *gin.Context, imageID int) (responses.ImageProduct, custom_error.CustomError)
	RemoveImage(c *gin.Context, imageID int) custom_error.CustomError
	UploadImage(c *gin.Context, productID int, imageStoredPath []string) custom_error.CustomError
}
