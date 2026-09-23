package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/domain/domain_err"
	"shop/interfaces/http/requests/web"
	"shop/interfaces/http/response"
	"shop/pkg/errors"
	"shop/pkg/old"
	"shop/pkg/sessions"
	"shop/pkg/util"
)

func (p PublicHandler) ShowProfile(c *gin.Context) {
	response.CustomerRender(c, http.StatusFound, "customer_profile", gin.H{
		"TITLE":  "مدیریت پروفایل",
		"ACTIVE": "profile",
	})
	return
}

func (p PublicHandler) EditProfile(c *gin.Context) {
	response.CustomerRender(c, http.StatusFound, "customer_edit_profile",
		gin.H{
			"TITLE":  "ویرایش پروفایل",
			"ACTIVE": "edit_profile",
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
		errors.SetErrors(c, p.dep.I18nBundle, bErr)
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	uErr := p.homeSrv.UpdateProfile(c, &req)
	if uErr.Code > 0 {
		sessions.Set(c, "message", domain_err.UpdateWasFailed)
	}
	sessions.Set(c, "message", domain_err.SuccessfullyUpdated)
	c.Redirect(http.StatusFound, "/profile/edit")
	c.Abort()
	return
}

func (p PublicHandler) StoreAddress(c *gin.Context) {
	var req requests.StoreAddressRequest
	_ = c.Request.ParseForm()
	err := c.ShouldBind(&req)
	if err != nil {

		errors.Init()
		errors.SetErrors(c, p.dep.I18nBundle, err)
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	//todo:add validation for postal_code
	if !util.ValidateIRMobile(req.ReceiverMobile) {
		errors.Init()

		errors.Add("receivermobile", domain_err.IRMobileIsInvalid)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	p.homeSrv.StoreAddress(c, &req)

	c.Redirect(http.StatusFound, c.Request.Referer())
	return
}
