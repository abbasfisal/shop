package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"shop/internal/modules/public/requests"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/custom_messages"
	"shop/internal/pkg/errors"
	"shop/internal/pkg/html"
	"shop/internal/pkg/sessions"
	"shop/internal/pkg/sms/kavenegar"
	"shop/internal/pkg/util"
	"time"
)

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

func (p PublicHandler) PostLogin(c *gin.Context) {
	var req requests.CustomerLoginRequest
	fmt.Println("--- step 1 ----")
	//bind
	_ = c.Request.ParseForm()
	if err := c.ShouldBind(&req); err != nil {
		errors.SetErrors(c, p.dep.I18nBundle, err)

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
	go kavenegar.SendOTP(req.Mobile, newOTP.Code)

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
	otpVerifyErr := p.homeSrv.VerifyOtp(c, mobile, &req)
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
