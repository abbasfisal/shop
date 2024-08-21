package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"shop/internal/middlewares"
	PublicHandler "shop/internal/modules/public/handlers"
	CustomerMiddlewares "shop/internal/modules/public/middlewares"
	"shop/internal/modules/public/services/home"
)

func SetPublic(r *gin.Engine, i18nBundle *i18n.Bundle) {

	homeSrv := home.NewHomeService()
	publicHdl := PublicHandler.NewPublicHandler(homeSrv, i18nBundle)

	r.GET("/", publicHdl.HomePage)
	r.GET("/products/single", publicHdl.SingleProduct)

	//find by sku (detail of  a product )
	r.GET("/:category_slug/:product_slug/:sku", publicHdl.ShowProduct)

	r.GET("/:category_slug", publicHdl.ShowProductsByCategory)

	publicAuthGrp := r.RouterGroup
	customerRoute := r.RouterGroup
	customerRoute.Use(CustomerMiddlewares.CheckCustomerSessionID())
	{
		customerRoute.GET("/login", publicHdl.ShowLogin)
		customerRoute.POST("/login", publicHdl.PostLogin)
		customerRoute.GET("/verify", publicHdl.ShowVerifyOtp)
		customerRoute.POST("/verify", publicHdl.PostVerifyOtp)
		customerRoute.GET("/resend-otp", publicHdl.ResendOtp)
	}

	publicAuthGrp.Use(CustomerMiddlewares.CustomerMustLogin())
	{
		publicAuthGrp.GET("/logout", publicHdl.LogOut)
		publicAuthGrp.GET("/profile", publicHdl.ShowProfile)
		publicAuthGrp.GET("/profile/edit", publicHdl.EditProfile)
		publicAuthGrp.POST("/profile/edit", publicHdl.UpdateProfile)
	}

	guestGrp := r.Group("/")
	guestGrp.Use(middlewares.IsGuest)
	{
		//guestGrp.GET("/users/login", publicHlr.ShowLogin)
		//guestGrp.POST("/users/login", publicHlr.PostLogin)
		//register routes
	}

	authGrp := r.Group("/")
	authGrp.Use(middlewares.IsAdmin)
	{

	}

}
