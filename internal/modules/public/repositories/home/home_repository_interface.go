package home

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
	"shop/internal/modules/public/requests"
	"shop/internal/pkg/custom_error"
)

type HomeRepositoryInterface interface {
	GetRandomProducts(ctx context.Context, limit int) ([]entities.Product, error)
	GetLatestProducts(ctx context.Context, limit int) ([]entities.Product, error)
	GetCategories(ctx context.Context, limit int) ([]entities.Category, error)
	GetProduct(ctx context.Context, productSlug, sku string) (entities.Product, error)
	GetProductsBy(ctx context.Context, columnName string, value any) ([]entities.Product, error)
	GetCategoryBy(ctx context.Context, columnName string, value any) (entities.Category, error)
	NewOtp(ctx context.Context, mobile string) (entities.OTP, custom_error.CustomError)
	VerifyOtp(c *gin.Context, mobile string, req requests.CustomerVerifyRequest) (entities.OTP, error)
	ProcessCustomerAuthenticate(c *gin.Context, mobile string) (entities.Session, error)
	LogOut(c *gin.Context) error
}
