package home

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"shop/internal/entities"
	"shop/internal/modules/admin/responses"
	"shop/internal/modules/public/repositories/home"
	"shop/internal/modules/public/repositories/home_mongo"
	"shop/internal/modules/public/requests"
	CustomerRes "shop/internal/modules/public/responses"
	"shop/internal/pkg/cache"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/custom_messages"
	"shop/internal/pkg/helpers"
	"shop/internal/pkg/pagination"
	"shop/internal/pkg/payment/zarinpal"
	"shop/internal/pkg/sessions"
)

type HomeService struct {
	repo      home.HomeRepositoryInterface
	mongoRepo home_mongo.MongoHomeRepositoryInterface
}

func NewHomeService(repo home.HomeRepositoryInterface, mongoRepo home_mongo.MongoHomeRepositoryInterface) HomeServiceInterface {
	return &HomeService{
		repo:      repo,
		mongoRepo: mongoRepo,
	}
}

//-----------------------------------
//<<<<<<<<<<<< Method >>>>>>>>>>>>>>>
//-----------------------------------

func (h *HomeService) GetProducts(ctx context.Context, limit int) (*responses.Products, custom_error.CustomError) {

	products, err := h.repo.GetLatestProducts(ctx, limit)
	if err != nil {
		return nil, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToProducts(products), custom_error.CustomError{}
}

func (h *HomeService) GetCategories(ctx context.Context, limit int) (*responses.Categories, custom_error.CustomError) {

	categories, err := h.repo.GetCategories(ctx, limit)
	if err != nil {
		return nil, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToCategories(categories), custom_error.CustomError{}
}

func (h *HomeService) ListProductByCategorySlug(c *gin.Context, slug string) (pagination.Pagination, error) {

	productList, err := h.repo.ListProductBy(c, slug)
	if err != nil {
		return pagination.Pagination{}, err
	}

	return productList, nil

}

func (h *HomeService) ShowCategory(ctx context.Context, columnName string, value any) (*responses.Category, custom_error.CustomError) {

	category, err := h.repo.GetCategoryBy(ctx, columnName, value)
	if err != nil {
		return nil, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToCategory(category), custom_error.CustomError{}
}

func (h *HomeService) SendOtp(ctx context.Context, Mobile string) (*entities.OTP, custom_error.CustomError) {
	return h.repo.NewOtp(ctx, Mobile)
}

func (h *HomeService) VerifyOtp(c *gin.Context, mobile string, req *requests.CustomerVerifyRequest) custom_error.CustomError {
	otp, err := h.repo.VerifyOtp(c, mobile, req)

	if err != nil {
		return custom_error.HandleError(err, custom_messages.OTPIsNotValid)
	}

	fmt.Println("--- home_service:VerifyOtp ---  otp is :--- ", otp)
	return custom_error.CustomError{}
}

func (h *HomeService) ProcessCustomerAuthentication(c *gin.Context, mobile string) (CustomerRes.CustomerSession, custom_error.CustomError) {
	sess, err := h.repo.ProcessCustomerAuthenticate(c, mobile)
	if err != nil {
		fmt.Println("------ error ProcessCustomerAuthentication: line : 99 ", err)
		return CustomerRes.CustomerSession{}, custom_error.New(err.Error(), "مشکل در ایجاد سشن", custom_error.CreateSessionFailedCode)
	}

	return CustomerRes.ToCustomerSession(sess), custom_error.CustomError{}
}

func (h *HomeService) LogOut(c *gin.Context) bool {
	err := h.repo.LogOut(c)
	if err != nil {
		return false
	}

	sessions.ClearAll(c)
	return true
}

func (h *HomeService) UpdateProfile(c *gin.Context, req *requests.CustomerProfileRequest) custom_error.CustomError {
	if err := h.repo.UpdateProfile(c, req); err != nil {
		fmt.Println("--- update profile failed : --- ", err)
		return custom_error.New(err.Error(), custom_error.SomethingWrongHappened, 500)
	}
	return custom_error.CustomError{}
}

func (h *HomeService) GetMenu(c context.Context) ([]*CustomerRes.CategoryResponse, error) {

	//get menu from cache
	menu := cache.Get(c, "menu")

	var categoryResponses []*CustomerRes.CategoryResponse

	if menu == "" {
		fmt.Println("--- menu was not exist in cache ------")

		//get menu from database
		menu, err := h.repo.GetMenu(c)
		if err != nil {
			return nil, err
		}

		for _, category := range menu {
			categoryResponse := CustomerRes.ToMenuResponse(category)
			categoryResponses = append(categoryResponses, categoryResponse)
		}

		//marsh repository response
		categoryJsonResponse, err := json.Marshal(categoryResponses)
		if err != nil {
			fmt.Println("--- category marshal error :", categoryJsonResponse)
		} else {
			fmt.Println("--- category marshal success :", categoryJsonResponse)
		}

		//store marshaled data into cache
		cacheSetErr := cache.Set(c, "menu", string(categoryJsonResponse), -1)
		if err != nil {
			fmt.Println("---- cache set menu key error: ", cacheSetErr)
		}

	} else {

		fmt.Println("--- menu was exist in cache ------")
		//menu was existed in cache
		unmarshalErr := json.Unmarshal([]byte(menu), &categoryResponses)
		if unmarshalErr != nil {
			fmt.Println("---- unmarshal category response err :", unmarshalErr)
		}
	}
	return categoryResponses, nil
}

func (h *HomeService) GetSingleProduct(c *gin.Context, productSku string, productSlug string) (map[string]interface{}, custom_error.CustomError) {

	mongoProduct, err := h.mongoRepo.GetProduct(c, productSku, productSlug)
	if err == nil {
		fmt.Println("--- mongo success --- ")
		return mongoProduct, custom_error.CustomError{}
	} else {
		fmt.Println("--- mongo product get err: ", err)
	}

	product, err := h.repo.GetProduct(c, productSku, productSlug)
	if err != nil {
		return nil, custom_error.HandleError(err, custom_error.RecordNotFound)
	}

	p := responses.ToProduct(product["product"].(*entities.Product))
	product["product"] = p

	return product, custom_error.CustomError{}
}

func (h *HomeService) AddToCart(c *gin.Context, productObjectID primitive.ObjectID, req requests.AddToCartRequest) {
	mongoProduct, err := h.mongoRepo.GetProductByObjectID(c, productObjectID, req)

	if err != nil {
		fmt.Println("error not found doc ")
		return
	}

	//store in cart
	user, ok := helpers.GetAuthUser(c)
	if !ok || user.ID <= 0 {
		return
	}

	h.repo.InsertCart(c, user, mongoProduct, req)

	fmt.Println("succ find :title", mongoProduct.ID)
}

func (h *HomeService) CartItemIncrement(c *gin.Context, req *requests.IncreaseCartItemQty) bool {
	err := h.repo.IncreaseCartItemCount(c, req)
	if err != nil {
		fmt.Println("[failed]-[CartItemIncrement]-[error]:", err)
		return false
	}
	return true
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
func (h *HomeService) ProcessOrderPayment(c *gin.Context, zarin *zarinpal.Zarinpal) (*entities.Order, *entities.Payment, custom_error.CustomError) {

	order, err := h.repo.GenerateOrderFromCart(c)
	if err != nil {
		return nil, nil, custom_error.New(err.Error(), custom_error.SomethingWrongHappened, 1000)
	}

	customer, _ := helpers.GetAuthUser(c)
	description := "order id :" + order.OrderNumber

	//paymentURL, authority, statusCode, zarinErr := zarin.NewPaymentRequest(int(order.TotalSalePrice), "http://vivify.ir/checkout/payment/verify", description, "", customer.Mobile)
	paymentURL, authority, statusCode, zarinErr := zarin.NewPaymentRequest(int(order.TotalSalePrice), os.Getenv("ZARINPAL_CALLBACKURL"), description, "", customer.Mobile)
	if zarinErr != nil || statusCode != 100 {
		log.Println("[home_service]-[ProcessOrderPayment]-[New ZarinPal Payment Request Error]:", zarinErr)
		return nil, nil, custom_error.New(err.Error(), custom_error.SomethingWrongHappened, 10001)
	}
	log.Println("[ZarinPal New Request Success]:", "paymentURL:", paymentURL, "|authority:", authority, "|statusCode:", statusCode)

	//create new payment
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
	paymentErr := h.repo.CreatePayment(c, &payment)
	if paymentErr != nil {
		return order, nil, custom_error.New(err.Error(), custom_error.SomethingWrongHappened, 10002)
	}

	return order, &payment, custom_error.CustomError{}
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
