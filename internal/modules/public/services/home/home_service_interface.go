package home

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
	"shop/internal/modules/admin/responses"
	"shop/internal/modules/public/requests"
	CustomerResp "shop/internal/modules/public/responses"
	"shop/internal/pkg/custom_error"
)

type HomeServiceInterface interface {
	GetProducts(ctx context.Context, limit int) (responses.Products, custom_error.CustomError)
	GetCategories(ctx context.Context, limit int) (responses.Categories, custom_error.CustomError)
	ShowCategory(ctx context.Context, columnName string, value any) (responses.Category, custom_error.CustomError)
	ShowProductDetail(ctx context.Context, productSlug, sku string) (responses.Product, custom_error.CustomError)
	ShowProductsByCategorySlug(ctx context.Context, value any) (responses.Products, custom_error.CustomError)

	// otp

	SendOtp(ctx context.Context, Mobile string) (entities.OTP, custom_error.CustomError)

	VerifyOtp(c *gin.Context, mobile string, req requests.CustomerVerifyRequest) custom_error.CustomError
	ProcessCustomerAuthentication(c *gin.Context, mobile string) (CustomerResp.CustomerSession, custom_error.CustomError)
}
