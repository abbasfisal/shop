package product

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/interfaces/http/requests/admin"
)

type ProductServiceInterface interface {
	Index(ctx context.Context) (*responses.Products, domain_err.CustomError)
	Show(ctx context.Context, columnName string, value any) (*responses.Product, []map[string]interface{}, domain_err.CustomError)
	Create(ctx context.Context, req *requests.CreateProductRequest) (*responses.Product, domain_err.CustomError)
	CheckSkuIsUnique(ctx context.Context, sku string) (bool, domain_err.CustomError)
	FetchByProductID(c *gin.Context, productID int) (*responses.Product, domain_err.CustomError)
	FetchRootAttributes(c *gin.Context, productID int) (*responses.Attributes, domain_err.CustomError)
	AddAttributeValues(c *gin.Context, productID int, attributes []string) domain_err.CustomError
	FetchProductAttributes(c *gin.Context, productID int) (map[string]interface{}, domain_err.CustomError)
	CreateInventory(c *gin.Context, productID int, req *requests.CreateProductInventoryRequest) domain_err.CustomError
	FetchImage(c *gin.Context, imageID int) (*responses.ImageProduct, domain_err.CustomError)
	RemoveImage(c *gin.Context, imageID int) domain_err.CustomError
	UploadImage(c *gin.Context, productID int, imageStoredPath []string) domain_err.CustomError
	Update(c *gin.Context, productID int, req *requests.UpdateProductRequest) domain_err.CustomError
	DeleteInventoryAttribute(c *gin.Context, productInventoryAttributeID int) domain_err.CustomError
	DeleteInventory(c *gin.Context, inventoryID int) domain_err.CustomError
	AppendAttributesToInventory(c *gin.Context, inventoryID int, attributes []string) domain_err.CustomError
	UpdateInventoryQuantity(c *gin.Context, inventoryID int, quantity uint) domain_err.CustomError
	AddFeature(c *gin.Context, productID int, req *requests.CreateProductFeatureRequest) domain_err.CustomError
	RemoveFeature(c *gin.Context, productID int, featureID int) domain_err.CustomError
	FetchFeature(c *gin.Context, productID int, featureID int) (*responses.Feature, domain_err.CustomError)
	UpdateFeature(c *gin.Context, productID int, featureID int, req *requests.UpdateProductFeatureRequest) domain_err.CustomError
	AddRecommendation(c *gin.Context, productID int, productRecommendationIDs []string) domain_err.CustomError
	FetchAllRecommendation(c *gin.Context, productID int) ([]map[string]interface{}, domain_err.CustomError)
}
