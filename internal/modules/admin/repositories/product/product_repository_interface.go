package product

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
)

type ProductRepositoryInterface interface {
	GetAll(ctx context.Context) ([]entities.Product, error)
	FindBy(ctx context.Context, columnName string, value any) (entities.Product, error)
	FindByID(ctx context.Context, ID int) (entities.Product, error)
	Store(ctx context.Context, product entities.Product) (entities.Product, error)
	GetRootAttributes(ctx *gin.Context, productID int) ([]entities.Attribute, error)
	StoreAttributeValues(ctx *gin.Context, productID int, attValues []string) error
}
