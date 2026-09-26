package repositories

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/domain/entities"
	"shop/interfaces/http/requests/admin"
)

type ProductRepositoryInterface interface {
	GetAll(ctx context.Context) ([]*entities.Product, error)
	// GetList is the admin list query: search + filters (status, category,
	// in_stock, attributes_json facet) + sorting (Laravel index()).
	GetList(ctx context.Context, q requests.ProductListQuery) ([]*entities.Product, error)
	FindBy(ctx context.Context, columnName string, value any) (*entities.Product, error)
	FindByID(ctx context.Context, ID int) (*entities.Product, error)
	// Store creates the product and its variants (rows) in one transaction.
	Store(ctx context.Context, product *entities.Product, rows []requests.VariantRow) (*entities.Product, error)
	// SyncReadModel rebuilds products.read_model (JSONB) and syncs Typesense.
	SyncReadModel(ctx context.Context, productID uint) error
	GetRootAttributes(ctx *gin.Context, productID int) ([]*entities.Attribute, error)
	StoreAttributeValues(ctx *gin.Context, productID int, attValues []string) error
	// DeleteProductAttribute removes one legacy product_attributes row.
	DeleteProductAttribute(c *gin.Context, productAttributeID int) (uint, error)
	GetProductAndAttributes(ctx *gin.Context, productID int) (map[string]interface{}, error)
	StoreProductInventory(c *gin.Context, productID int, req *requests.CreateProductInventoryRequest) (*entities.ProductVariant, error)
	GetImage(c *gin.Context, imageID int) (*entities.ProductImages, error)
	DeleteImage(c *gin.Context, imageID int) error
	StoreImages(c *gin.Context, productID int, imageStoredPath []string) error
	// Update writes the product columns and its variant rows in one
	// transaction (rows with ID > 0 update, the rest are created).
	Update(c *gin.Context, productID int, req *requests.UpdateProductRequest) (*entities.Product, error)
	// inventory/variant mutations return the affected product id so the
	// application layer can refresh pricing aggregates afterwards
	DeleteInventoryAttribute(c *gin.Context, variantAttributeValueID int) (uint, error)
	DeleteInventory(c *gin.Context, variantID int) (uint, error)
	AppendAttributesToInventory(c *gin.Context, variantID int, attributes []string) (uint, error)
	UpdateInventoryQuantity(c *gin.Context, variantID int, quantity uint) (uint, error)
	InsertFeature(c *gin.Context, productID int, req *requests.CreateProductFeatureRequest) error
	DeleteFeature(c *gin.Context, productID int, featureID int) error
	GetFeatureBy(c *gin.Context, productID int, featureID int) (*entities.Feature, error)
	EditFeature(c *gin.Context, productID int, featureID int, req *requests.UpdateProductFeatureRequest) error
	// GetAllProductBriefs returns lightweight product summaries for the admin
	// recommendation picker (shape consumed by admin_edit_product.html).
	GetAllProductBriefs(c context.Context) ([]map[string]interface{}, error)
	InsertRecommendation(c *gin.Context, productID int, productRecommendationIDs []string) error
	GetAllRecommendation(c *gin.Context, productID int) ([]map[string]interface{}, error)
	// RefreshProductAggregates recomputes min/max price, stock totals,
	// in_stock, variants_count, product_type and attributes_json (JSONB).
	RefreshProductAggregates(ctx context.Context, productID uint) error
}
