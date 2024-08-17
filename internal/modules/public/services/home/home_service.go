package home

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
	"shop/internal/modules/admin/responses"
	"shop/internal/modules/public/repositories/home"
	"shop/internal/modules/public/requests"
	CustomerRes "shop/internal/modules/public/responses"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/custom_messages"
)

type HomeService struct {
	repo home.HomeRepositoryInterface
}

func NewHomeService() HomeService {
	return HomeService{
		repo: home.NewHomeRepository(),
	}
}

func (h HomeService) GetProducts(ctx context.Context, limit int) (responses.Products, custom_error.CustomError) {

	products, err := h.repo.GetLatestProducts(ctx, limit)
	if err != nil {
		return responses.Products{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToProducts(products), custom_error.CustomError{}
}

func (h HomeService) GetCategories(ctx context.Context, limit int) (responses.Categories, custom_error.CustomError) {

	categories, err := h.repo.GetCategories(ctx, limit)
	if err != nil {
		return responses.Categories{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToCategories(categories), custom_error.CustomError{}
}

func (h HomeService) ShowProductDetail(ctx context.Context, productSlug, sku string) (responses.Product, custom_error.CustomError) {

	product, err := h.repo.GetProduct(ctx, productSlug, sku)
	if err != nil {
		return responses.Product{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToProduct(product), custom_error.CustomError{}
}

func (h HomeService) ShowProductsByCategorySlug(ctx context.Context, value any) (responses.Products, custom_error.CustomError) {

	//get
	category, cErr := h.ShowCategory(ctx, "slug", value)
	if cErr.Code == 404 {
		return responses.Products{}, custom_error.New(cErr.Error(), custom_error.RecordNotFound, 404)
	}
	if cErr.Code == 500 {
		return responses.Products{}, custom_error.New(cErr.Error(), custom_error.InternalServerError, 500)
	}

	products, err := h.repo.GetProductsBy(ctx, "category_id", category.ID)
	if err != nil {
		return responses.Products{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToProducts(products), custom_error.CustomError{}
}

func (h HomeService) ShowCategory(ctx context.Context, columnName string, value any) (responses.Category, custom_error.CustomError) {

	category, err := h.repo.GetCategoryBy(ctx, columnName, value)
	if err != nil {
		return responses.Category{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToCategory(category), custom_error.CustomError{}
}

//

func (h HomeService) SendOtp(ctx context.Context, Mobile string) (entities.OTP, custom_error.CustomError) {
	return h.repo.NewOtp(ctx, Mobile)
}

func (h HomeService) VerifyOtp(c *gin.Context, mobile string, req requests.CustomerVerifyRequest) custom_error.CustomError {
	otp, err := h.repo.VerifyOtp(c, mobile, req)

	if err != nil {
		return custom_error.HandleError(err, custom_messages.OTPIsNotValid)
	}

	fmt.Println("--- home_service:VerifyOtp ---  otp is :--- ", otp)
	return custom_error.CustomError{}
}

func (h HomeService) ProcessCustomerAuthentication(c *gin.Context, mobile string) (CustomerRes.CustomerSession, custom_error.CustomError) {
	sess, err := h.repo.ProcessCustomerAuthenticate(c, mobile)
	if err != nil {
		fmt.Println("------ error ProcessCustomerAuthentication: line : 99 ", err)
		return CustomerRes.CustomerSession{}, custom_error.New(err.Error(), "مشکل در ایجاد سشن", custom_error.CreateSessionFailedCode)
	}

	return CustomerRes.ToCustomerSession(sess), custom_error.CustomError{}
}
