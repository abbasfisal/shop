package home

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"shop/application/dto/admin"
	CustomerRes "shop/application/dto/web"
	"shop/bootstrap"
	"shop/domain/domain_err"
	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/infrastructure/events"
	"shop/infrastructure/messages"
	"shop/infrastructure/payment/zarinpal"
	"shop/interfaces/http/requests/web"
	"shop/pkg/cache"
	"shop/pkg/helpers"
	"shop/pkg/pagination"
	"shop/pkg/sessions"
	"time"
)

type HomeService struct {
	dep  *bootstrap.Dependencies
	repo repositories.HomeRepositoryInterface
}

func NewHomeService(dep *bootstrap.Dependencies, repo repositories.HomeRepositoryInterface, eventManager *events.EventManager) HomeServiceInterface {

	return &HomeService{
		dep:  dep,
		repo: repo,
	}
}

//-----------------------------------
//<<<<<<<<<<<< Method >>>>>>>>>>>>>>>
//-----------------------------------

func (h *HomeService) GetProducts(ctx context.Context, limit int) (*responses.Products, domain_err.CustomError) {

	products, err := h.repo.GetLatestProducts(ctx, limit)
	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToProducts(products), domain_err.CustomError{}
}

func (h *HomeService) GetCategories(ctx context.Context, limit int) (*responses.Categories, domain_err.CustomError) {

	categories, err := h.repo.GetCategories(ctx, limit)
	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToCategories(categories), domain_err.CustomError{}
}

func (h *HomeService) ListProductByCategorySlug(c *gin.Context, slug string) (pagination.Pagination, error) {

	productList, err := h.repo.ListProductBy(c, slug)
	if err != nil {
		return pagination.Pagination{}, err
	}

	return productList, nil

}

func (h *HomeService) ShowCategory(ctx context.Context, columnName string, value any) (*responses.Category, domain_err.CustomError) {

	category, err := h.repo.GetCategoryBy(ctx, columnName, value)
	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToCategory(category), domain_err.CustomError{}
}

func (h *HomeService) SendOtp(ctx context.Context, Mobile string) (*entities.OTP, domain_err.CustomError) {
	return h.repo.NewOtp(ctx, Mobile)
}

func (h *HomeService) VerifyOtp(c *gin.Context, mobile string, req *requests.CustomerVerifyRequest) domain_err.CustomError {
	otp, err := h.repo.VerifyOtp(c, mobile, req)

	if err != nil {
		return domain_err.HandleError(err, custom_messages.OTPIsNotValid)
	}

	fmt.Println("--- home_service:VerifyOtp ---  otp is :--- ", otp)
	return domain_err.CustomError{}
}

func (h *HomeService) ProcessCustomerAuthentication(c *gin.Context, mobile string) (CustomerRes.CustomerSession, domain_err.CustomError) {
	sess, err := h.repo.ProcessCustomerAuthenticate(c, mobile)
	if err != nil {
		fmt.Println("------ error ProcessCustomerAuthentication: line : 99 ", err)
		return CustomerRes.CustomerSession{}, domain_err.New(err.Error(), "مشکل در ایجاد سشن", domain_err.CreateSessionFailedCode)
	}

	return CustomerRes.ToCustomerSession(sess), domain_err.CustomError{}
}

func (h *HomeService) LogOut(c *gin.Context) bool {
	err := h.repo.LogOut(c)
	if err != nil {
		return false
	}

	sessions.ClearAll(c)
	return true
}

func (h *HomeService) UpdateProfile(c *gin.Context, req *requests.CustomerProfileRequest) domain_err.CustomError {
	if err := h.repo.UpdateProfile(c, req); err != nil {
		fmt.Println("--- update profile failed : --- ", err)
		return domain_err.New(err.Error(), domain_err.SomethingWrongHappened, 500)
	}
	return domain_err.CustomError{}
}

// menuCacheTTL keeps the header menu fresh even when a category row is
// written outside the repository (e.g. by the seeder, which cannot invalidate
// the cache). Category writes still invalidate it immediately.
const menuCacheTTL = 5 * time.Minute

func (h *HomeService) GetMenu(c context.Context) ([]*CustomerRes.CategoryResponse, error) {
	if menu, ok := menuFromCache(c); ok {
		return menu, nil
	}

	// cache miss / stale entry → read the categories from the database
	menu, err := h.repo.GetMenu(c)
	if err != nil {
		return nil, err
	}

	categoryResponses := make([]*CustomerRes.CategoryResponse, 0, len(menu))
	for _, category := range menu {
		categoryResponses = append(categoryResponses, CustomerRes.ToMenuResponse(category))
	}

	// never cache an empty menu: json.Marshal(nil slice) is the literal
	// "null", and with a no-expiry key that value poisoned the menu forever
	// (the header showed an empty «دسته بندی کالاها» panel).
	if len(categoryResponses) > 0 {
		payload, marshalErr := json.Marshal(categoryResponses)
		if marshalErr != nil {
			fmt.Println("---- menu marshal error :", marshalErr)
		} else if setErr := cache.Set(c, "menu", string(payload), menuCacheTTL); setErr != nil {
			fmt.Println("---- cache set menu key error: ", setErr)
		}
	}
	return categoryResponses, nil
}

// menuFromCache returns the cached menu only when it is a usable payload.
// "null" / "[]" / unparsable / empty entries are treated as a miss so the
// next request rebuilds them from the database.
func menuFromCache(c context.Context) ([]*CustomerRes.CategoryResponse, bool) {
	raw := cache.Get(c, "menu")
	if raw == "" || raw == "null" || raw == "[]" {
		return nil, false
	}

	var menu []*CustomerRes.CategoryResponse
	if err := json.Unmarshal([]byte(raw), &menu); err != nil || len(menu) == 0 {
		return nil, false
	}
	return menu, true
}

func (h *HomeService) GetSingleProduct(c *gin.Context, productSku string, productSlug string) (map[string]interface{}, []entities.RecommendedProduct, domain_err.CustomError) {

	product, recommendations, err := h.repo.GetProduct(c, productSku, productSlug)
	if err != nil {
		return nil, nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}

	return product, recommendations, domain_err.CustomError{}
}

func (h *HomeService) AddToCart(c *gin.Context, productID uint, req requests.AddToCartRequest) {
	prod, err := h.repo.GetProductByID(c, productID)
	if err != nil {
		fmt.Println("[error]-[AddToCart]: product not found, id:", productID)
		return
	}

	//store in cart
	user, ok := helpers.GetAuthUser(c)
	if !ok || user.ID <= 0 {
		return
	}

	h.repo.InsertCart(c, user, prod, req)

	fmt.Println("succ find :title", prod.Title)
}

func (h *HomeService) CartItemIncrement(c *gin.Context, req *requests.IncreaseCartItemQty) error {
	err := h.repo.IncreaseCartItemCount(c, req)

	if err != nil {
		return err
	}
	return nil
}

func (h *HomeService) CartItemDecrement(c *gin.Context, req *requests.IncreaseCartItemQty) bool {
	err := h.repo.DecreaseCartItemCount(c, req)
	if err != nil {
		fmt.Println("[failed]-[CartItemDecrement]-[error]:", err)
		return false
	}
	return true
}

func (h *HomeService) RemoveCartItem(c *gin.Context, req *requests.IncreaseCartItemQty) bool {
	err := h.repo.DeleteCartItem(c, req)
	if err != nil {
		fmt.Println("[failed]-[RemoveCartItem]-[error]:", err)
		return false
	}

	return true
}

func (h *HomeService) StoreAddress(c *gin.Context, req *requests.StoreAddressRequest) {
	err := h.repo.CreateOrUpdateAddress(c, req)
	if err != nil {
		fmt.Println("[home_service]-[StoreAddress]-err:", err)
	}
}

// ProcessOrderPayment convert cart to order and remove cart
func (h *HomeService) ProcessOrderPayment(c *gin.Context, zarin *zarinpal.Zarinpal) (*entities.Order, *entities.Payment, uint, error) {

	t := time.Now()
	order, inventoryID, err := h.repo.GenerateOrderFromCart(c)
	s := time.Since(t)
	fmt.Println("time left : ", s)

	if err != nil {
		if errors.Is(err, domain_err.OutOfStock) {
			return nil, nil, inventoryID, domain_err.OutOfStock
		}
		return nil, nil, 0, domain_err.InternalServerErr
	}

	customer, _ := helpers.GetAuthUser(c)
	description := "order id :" + order.OrderNumber

	//paymentURL, authority, statusCode, zarinErr := zarin.NewPaymentRequest(int(order.TotalSalePrice), "http://vivify.ir/checkout/payment/verify", description, "", customer.Mobile)
	paymentURL, authority, statusCode, zarinErr := zarin.NewPaymentRequest(int(order.TotalSalePrice), os.Getenv("ZARINPAL_CALLBACKURL"), description, "", customer.Mobile)
	if zarinErr != nil || statusCode != 100 {
		log.Println("[home_service]-[ProcessOrderPayment]-[New ZarinPal Payment Request Error]:", zarinErr)
		return nil, nil, 0, domain_err.InternalServerErr
	}
	log.Println("[ZarinPal New Request Success]:", "paymentURL:", paymentURL, "|authority:", authority, "|statusCode:", statusCode)

	//prepare new payment
	payment := entities.Payment{
		CustomerID:  customer.ID,
		OrderID:     order.ID,
		Authority:   authority,
		Description: description,
		PaymentURL:  paymentURL,
		StatusCode:  statusCode,
		Amount:      order.TotalSalePrice,
		RefID:       "",
		Status:      0, //pending
	}

	// store new payment in db
	paymentErr := h.repo.CreatePayment(c, &payment)
	if paymentErr != nil {
		return order, nil, 0, paymentErr
	}

	return order, &payment, 0, nil
}

func (h *HomeService) VerifyPayment(c *gin.Context, order *entities.Order, refID string, verified bool) {
	h.repo.OrderPaidSuccessfully(c, order, refID, verified)
}

func (h *HomeService) GetPaymentBy(c *gin.Context, authority string) (*entities.Order, entities.Customer, error) {
	return h.repo.GetPayment(c, authority)
}

func (h *HomeService) ListOrders(c *gin.Context) (pagination.Pagination, error) {
	orderList, err := h.repo.GetPaginatedOrders(c)
	if err != nil {
		return pagination.Pagination{}, err
	}

	//orderList.Rows = responses.ToOrders(orderList.Rows.([]*entities.Order))

	return orderList, nil

}

func (h *HomeService) GetOrderBy(c *gin.Context, orderNumber string) (interface{}, interface{}) {
	order, err := h.repo.GetOrder(c, orderNumber)
	if order == nil || err != nil {
		return nil, err
	}
	return CustomerRes.ToCustomerOrder(order), err
}
