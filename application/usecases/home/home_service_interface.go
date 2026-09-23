package home

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/application/dto/admin"
	CustomerResp "shop/application/dto/web"
	"shop/domain/domain_err"
	"shop/domain/entities"
	"shop/infrastructure/payment/zarinpal"
	"shop/interfaces/http/requests/web"
	"shop/pkg/pagination"
)

type HomeServiceInterface interface {
	GetProducts(ctx context.Context, limit int) (*responses.Products, domain_err.CustomError)
	GetCategories(ctx context.Context, limit int) (*responses.Categories, domain_err.CustomError)
	ShowCategory(ctx context.Context, columnName string, value any) (*responses.Category, domain_err.CustomError)

	ListProductByCategorySlug(c *gin.Context, slug string) (pagination.Pagination, error)

	// GetMenu fetch categories to show in menu
	GetMenu(c context.Context) ([]*CustomerResp.CategoryResponse, error)
	SendOtp(ctx context.Context, Mobile string) (*entities.OTP, domain_err.CustomError)
	VerifyOtp(c *gin.Context, mobile string, req *requests.CustomerVerifyRequest) domain_err.CustomError
	ProcessCustomerAuthentication(c *gin.Context, mobile string) (CustomerResp.CustomerSession, domain_err.CustomError)
	LogOut(c *gin.Context) bool
	UpdateProfile(c *gin.Context, req *requests.CustomerProfileRequest) domain_err.CustomError
	GetSingleProduct(c *gin.Context, productSku string, productSlug string) (map[string]interface{}, []entities.RecommendedProduct, domain_err.CustomError)

	//------cart

	AddToCart(c *gin.Context, productID uint, req requests.AddToCartRequest)
	CartItemIncrement(c *gin.Context, req *requests.IncreaseCartItemQty) error
	CartItemDecrement(c *gin.Context, req *requests.IncreaseCartItemQty) bool
	RemoveCartItem(c *gin.Context, req *requests.IncreaseCartItemQty) bool
	StoreAddress(c *gin.Context, req *requests.StoreAddressRequest)

	// ProcessOrderPayment convert cart to order and remove cart
	ProcessOrderPayment(c *gin.Context, zarin *zarinpal.Zarinpal) (*entities.Order, *entities.Payment, uint, error)

	VerifyPayment(c *gin.Context, payment *entities.Order, refID string, verified bool)
	GetPaymentBy(c *gin.Context, authority string) (*entities.Order, entities.Customer, error)

	ListOrders(c *gin.Context) (pagination.Pagination, error)
	GetOrderBy(c *gin.Context, orderNumber string) (interface{}, interface{})
}
