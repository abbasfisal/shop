package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"net/http"
	"shop/internal/modules/public/requests"
	"shop/internal/modules/public/services/home"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/custom_messages"
	"shop/internal/pkg/errors"
	"shop/internal/pkg/html"
	"shop/internal/pkg/old"
	"shop/internal/pkg/sessions"
	"shop/internal/pkg/util"
	"time"
)

type PublicHandler struct {
	homeSrv    home.HomeServiceInterface
	i18nBundle *i18n.Bundle
}

func NewPublicHandler(homeSrv home.HomeServiceInterface, i18nBundle *i18n.Bundle) PublicHandler {
	return PublicHandler{
		homeSrv:    homeSrv,
		i18nBundle: i18nBundle,
	}
}

func (p PublicHandler) Index(c *gin.Context) {
	//20 latest
	//20 random /-
	//20 lowest quantity
	//10 category
	//filter by price
	//filter by quantity
	//pagination in selected category

	products, err := p.homeSrv.GetProducts(context.TODO(), 20)
	if err.Code == 404 {
		c.JSON(200, gin.H{"err": 404})
		return
	}
	if err.Code == 500 {
		c.JSON(200, gin.H{"err": 500})
		return
	}

	categories, cErr := p.homeSrv.GetCategories(context.TODO(), 20)
	if cErr.Code == 404 {
		c.JSON(200, gin.H{"err": 404})
		return
	}
	if cErr.Code == 500 {
		c.JSON(200, gin.H{"err": 500})
		return
	}

	c.JSON(200, gin.H{
		"PRODUCTS":   products,
		"CATEGORIES": categories,
	})
	//	html.Render(c, http.StatusFound, "home", gin.H{
	//			"TITLE":      "home page",
	//			"PRODUCTS":   "",
	//			"CATEGORIES": "",
	//		})
	return
}

func (p PublicHandler) ShowProduct(c *gin.Context) {

	_ = c.Param("category_slug")
	productSlug := c.Param("product_slug")
	sku := c.Param("sku")

	product, err := p.homeSrv.ShowProductDetail(c, productSlug, sku)
	if err.Code == 404 {
		c.Redirect(http.StatusFound, "/")
		return
	}
	if err.Code == 500 {
		html.Error500(c)
		return
	}

	html.Render(c, http.StatusFound, "product_detail", gin.H{
		"PRODUCT": product,
	})
	return
}

func (p PublicHandler) ShowProductsByCategory(c *gin.Context) {

	products, err := p.homeSrv.ShowProductsByCategorySlug(c, c.Param("category_slug"))
	if err.Code == 404 {
		c.JSON(200, gin.H{
			"msg": "not found",
		})
		return
	}
	if err.Code == 500 {
		c.JSON(200, gin.H{
			"msg": "internal server err",
		})
		return
	}

	html.Render(c, http.StatusFound, "products_by_category_slug", gin.H{
		"PRODUCTS": products,
	})
	return

}

func (p PublicHandler) ShowLogin(c *gin.Context) {

	if sessions.GET(c, "otp_created_at") != "" || sessions.GET(c, "mobile") != "" {
		c.Redirect(http.StatusFound, "/verify")
		return
	}

	html.CustomerRender(c, 200, "customer_login", gin.H{
		"TITLE": "اسم فروشگاه",
		"data":  "data",
	})
}

func (p PublicHandler) ShowVerifyOtp(c *gin.Context) {
	//sessions.Set(c, "mobile", req.Mobile)
	//	sessions.Set(c, "otp_to_expire", otpToExpire)
	//	sessions.Set(c, "otp_created_at", otpCreateAt)

	if sessions.GET(c, "otp_created_at") == "" || sessions.GET(c, "mobile") == "" {
		fmt.Println(" ==== show verify otp form null ----- ")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	otpCreatedAt, err := time.Parse(time.RFC3339, sessions.GET(c, "otp_created_at"))
	if err != nil {
		sessions.ClearAll(c)
		fmt.Println(" ==== show verify otp created at ----- ")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	fmt.Println("====== time since ========: ", time.Since(otpCreatedAt))

	otpTTL := time.Duration(viper.GetInt("app.otp_expiration_time")) * time.Minute //in minute
	if time.Since(otpCreatedAt) > otpTTL {
		sessions.ClearAll(c)
		fmt.Println(" ==== otp is expired ====== ")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	html.CustomerRender(c, 200, "customer_verify_phone_number", gin.H{
		"TITLE":    "خرید از باآف باکیفیت و مقرون به صرفه",
		"MOBILE":   sessions.GET(c, "mobile"),
		"TOEXPIRE": otpTTL.Seconds() - time.Since(otpCreatedAt).Seconds(),
	})
}

func (p PublicHandler) HomePage(c *gin.Context) {
	menu, err := p.homeSrv.GetMenu(c)
	if err != nil {
		html.Error500(c)
		return
	}
	//row-header
	//row-newest
	//row-random
	//row-banners
	//row-by-category
	//

	html.CustomerRender(c, 200, "home", gin.H{
		"TITLE": "صفحه اصلی فروشگاه",
		"MENU":  menu,
	})
}

func (p PublicHandler) SingleProduct(c *gin.Context) {
	html.Render(c, 200, "customer_single_product", gin.H{
		"TITLE": "عنوان فروشگاه",
	})
}

func (p PublicHandler) PostLogin(c *gin.Context) {
	var req requests.CustomerLoginRequest
	fmt.Println("--- step 1 ----")
	//bind
	_ = c.Request.ParseForm()
	if err := c.ShouldBind(&req); err != nil {
		errors.SetErrors(c, p.i18nBundle, err)

		c.Redirect(http.StatusFound, "/login")
		return
	}

	fmt.Println("--- step 2 ----")

	//validate mobile
	if !util.ValidateIRMobile(req.Mobile) {
		errors.Init()
		errors.Add("mobile", "شماره موبایل معتبر وارد کنید")
		sessions.Set(c, "errors", errors.ToString())
		c.Redirect(http.StatusFound, "/login")
		return
	}

	fmt.Println("--- step 3 ----")

	newOTP, otpErr := p.homeSrv.SendOtp(c, req.Mobile)
	if otpErr.Code > 0 {
		if otpErr.Code == custom_error.OTPTooSoonCode {
			sessions.Set(c, "message", custom_error.OTPRequestTooSoon)
			fmt.Println("------ redirect to verify : to soon request : ")
			c.Redirect(http.StatusFound, "/verify")
			return
		}
		sessions.Set(c, "message", otpErr.DisplayMessage)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	fmt.Println("--- step 4 ----")

	fmt.Println("-------- new otp generated----- : ", newOTP)

	otpCreateAt := time.Now().Format(time.RFC3339)

	sessions.Set(c, "mobile", req.Mobile)
	sessions.Set(c, "otp_created_at", otpCreateAt)

	fmt.Println("--- step 6 ----")

	c.Redirect(http.StatusFound, "/verify")
	return

}

func (p PublicHandler) PostVerifyOtp(c *gin.Context) {
	mobile := sessions.GET(c, "mobile")
	otpCreatedAt, err := time.Parse(time.RFC3339, sessions.GET(c, "otp_created_at"))

	if mobile == "" || sessions.GET(c, "otp_created_at") == "" || !util.ValidateIRMobile(mobile) || err != nil {
		sessions.ClearAll(c)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	//check otp is expired?
	otpTTL := time.Duration(viper.GetInt("app.otp_expiration_time")) * time.Minute //in minute
	if time.Since(otpCreatedAt) > otpTTL {
		sessions.Set(c, "message", custom_messages.OTPISExpired)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	//bind
	var req requests.CustomerVerifyRequest
	bindErr := c.ShouldBind(&req)
	if bindErr != nil {
		sessions.Set(c, "message", custom_error.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	//verify otp
	otpVerifyErr := p.homeSrv.VerifyOtp(c, mobile, req)
	if otpVerifyErr.Code == 404 {
		sessions.Set(c, "message", custom_messages.OTPIsNotValid)
		c.Redirect(http.StatusFound, "/verify")
		c.Abort()
		return
	}
	if otpVerifyErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	//create record in sessions and customers table
	customerSession, pcaErr := p.homeSrv.ProcessCustomerAuthentication(c, mobile)
	if pcaErr.Code > 0 {
		fmt.Println("---- handler PostVerifyOtp --- err : ", pcaErr.Error())
		sessions.Set(c, "message", custom_error.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/verify")
		return
	}

	//store session(uuid) in session
	sessions.Set(c, "session_id", customerSession.SessionID)
	c.Redirect(http.StatusFound, "/")
	return
}

func (p PublicHandler) ResendOtp(c *gin.Context) {
	mobile := sessions.GET(c, "mobile")
	otpCreatedAt, err := time.Parse(time.RFC3339, sessions.GET(c, "otp_created_at"))

	if mobile == "" || sessions.GET(c, "otp_created_at") == "" || !util.ValidateIRMobile(mobile) || err != nil {
		sessions.ClearAll(c)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	//check otp expire time
	otpTTL := time.Duration(viper.GetInt("app.otp_expiration_time")) * time.Minute //in minute
	if time.Since(otpCreatedAt) < otpTTL {
		c.Redirect(http.StatusFound, "/verify")
		return
	}

	//otp is expired , resend new otp
	newOTP, otpErr := p.homeSrv.SendOtp(c, mobile)
	if otpErr.Code > 0 {
		if otpErr.Code == custom_error.OTPTooSoonCode {
			sessions.Set(c, "message", custom_error.OTPRequestTooSoon)
			fmt.Println("------ redirect to verify : to soon request : ")
			c.Redirect(http.StatusFound, "/verify")
			return
		}
		sessions.Set(c, "message", otpErr.DisplayMessage)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	fmt.Println("-------- new otp generated----- : ", newOTP)
	otpCreateAt := time.Now().Format(time.RFC3339)

	fmt.Println("\n ----- otp created at ---- : ", otpCreateAt)

	sessions.Set(c, "mobile", mobile)
	sessions.Set(c, "otp_created_at", otpCreateAt)

	c.Redirect(http.StatusFound, "/verify")
	return

}

func (p PublicHandler) LogOut(c *gin.Context) {

	result := p.homeSrv.LogOut(c)

	if !result {
		fmt.Println("--- logOut was failed ;) ---- ")
	}
	fmt.Println("---- logOut was success ----")
	c.Redirect(http.StatusFound, "/")
	c.Abort()
	return
}

func (p PublicHandler) ShowProfile(c *gin.Context) {
	html.CustomerRender(c, http.StatusFound, "customer_profile", gin.H{
		"TITLE": "مدیریت پروفایل",
	})
	return
}

func (p PublicHandler) EditProfile(c *gin.Context) {
	html.CustomerRender(c, http.StatusFound, "customer_edit_profile",
		gin.H{
			"TITLE": "ویرایش پروفایل",
		},
	)
	return
}

func (p PublicHandler) UpdateProfile(c *gin.Context) {
	var req requests.CustomerProfileRequest

	//binding
	bErr := c.ShouldBind(&req)
	if bErr != nil {
		fmt.Println("-- err : -- ", bErr)
		errors.Init()
		errors.SetErrors(c, p.i18nBundle, bErr)
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	uErr := p.homeSrv.UpdateProfile(c, req)
	if uErr.Code > 0 {
		sessions.Set(c, "message", custom_error.UpdateWasFailed)
	}
	sessions.Set(c, "message", custom_error.SuccessfullyUpdated)
	c.Redirect(http.StatusFound, "/profile/edit")
	c.Abort()
	return
}
