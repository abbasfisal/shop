package product

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"shop/internal/entities"
	"shop/internal/modules/admin/requests"
)

type ProductRepositoryInterface interface {
	GetAll(ctx context.Context) ([]*entities.Product, error)
	FindBy(ctx context.Context, columnName string, value any) (*entities.Product, error)
	FindByID(ctx context.Context, ID int) (*entities.Product, error)
	Store(ctx context.Context, product *entities.Product) (*entities.Product, error)
	GetRootAttributes(ctx *gin.Context, productID int) ([]*entities.Attribute, error)
	StoreAttributeValues(ctx *gin.Context, productID int, attValues []string) error
	GetProductAndAttributes(ctx *gin.Context, productID int) (map[string]interface{}, error)
	StoreProductInventory(c *gin.Context, productID int, req *requests.CreateProductInventoryRequest) (*entities.ProductInventory, error)
	GetImage(c *gin.Context, imageID int) (*entities.ProductImages, error)
	DeleteImage(c *gin.Context, imageID int) error
	StoreImages(c *gin.Context, productID int, imageStoredPath []string) error
	Update(c *gin.Context, productID int, req *requests.UpdateProductRequest) (*entities.Product, error)
	DeleteInventoryAttribute(c *gin.Context, inventoryID int) error
	DeleteInventory(c *gin.Context, inventoryID int) error
	AppendAttributesToInventory(c *gin.Context, inventoryID int, attributes []string) error
	UpdateInventoryQuantity(c *gin.Context, inventoryID int, quantity uint) error
	InsertFeature(c *gin.Context, productID int, req *requests.CreateProductFeatureRequest) error
	DeleteFeature(c *gin.Context, productID int, featureID int) error
	GetFeatureBy(c *gin.Context, productID int, featureID int) (*entities.Feature, error)
	EditFeature(c *gin.Context, productID int, featureID int, req *requests.UpdateProductFeatureRequest) error
	GetAllMongoProduct(c context.Context) ([]bson.M, error)
	InsertRecommendation(c *gin.Context, productID int, productRecommendationIDs []string) error
	GetAllRecommendation(c *gin.Context, productID int) ([]bson.M, error)
}
