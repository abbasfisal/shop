package home

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
	"shop/internal/modules/admin/responses"
	"shop/internal/modules/public/requests"
	CustomerResp "shop/internal/modules/public/responses"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/pagination"
)

type HomeServiceInterface interface {
	GetProducts(ctx context.Context, limit int) (responses.Products, custom_error.CustomError)
	GetCategories(ctx context.Context, limit int) (responses.Categories, custom_error.CustomError)
	ShowCategory(ctx context.Context, columnName string, value any) (responses.Category, custom_error.CustomError)
	ShowProductDetail(ctx context.Context, productSlug, sku string) (responses.Product, custom_error.CustomError)
	ListProductByCategorySlug(c *gin.Context, slug string) (pagination.Pagination, error)

	// GetMenu fetch categories to show in menu
	GetMenu(c context.Context) ([]CustomerResp.CategoryResponse, error)

	// otp

	SendOtp(ctx context.Context, Mobile string) (entities.OTP, custom_error.CustomError)

	VerifyOtp(c *gin.Context, mobile string, req requests.CustomerVerifyRequest) custom_error.CustomError
	ProcessCustomerAuthentication(c *gin.Context, mobile string) (CustomerResp.CustomerSession, custom_error.CustomError)
	LogOut(c *gin.Context) bool
	UpdateProfile(c *gin.Context, req requests.CustomerProfileRequest) custom_error.CustomError
	GetSingleProduct(c *gin.Context, productSku string, productSlug string) (map[string]interface{}, custom_error.CustomError)
}
