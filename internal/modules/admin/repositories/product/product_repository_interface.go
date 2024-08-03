package product

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
	"shop/internal/modules/admin/requests"
)

type ProductRepositoryInterface interface {
	GetAll(ctx context.Context) ([]entities.Product, error)
	FindBy(ctx context.Context, columnName string, value any) (entities.Product, error)
	FindByID(ctx context.Context, ID int) (entities.Product, error)
	Store(ctx context.Context, product entities.Product) (entities.Product, error)
	GetRootAttributes(ctx *gin.Context, productID int) ([]entities.Attribute, error)
	StoreAttributeValues(ctx *gin.Context, productID int, attValues []string) error
	GetProductAndAttributes(ctx *gin.Context, productID int) (entities.Product, error)
	StoreProductInventory(c *gin.Context, productID int, req requests.CreateProductInventoryRequest) (entities.ProductInventory, error)
	GetImage(c *gin.Context, imageID int) (entities.ProductImages, error)
	DeleteImage(c *gin.Context, imageID int) error
	StoreImages(c *gin.Context, productID int, imageStoredPath []string) error
}
