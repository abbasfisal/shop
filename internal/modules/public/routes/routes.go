package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"shop/internal/middlewares"
	PublicHandler "shop/internal/modules/public/handlers"
	"shop/internal/modules/public/services/home"
)

func SetPublic(r *gin.Engine, i18nBundle *i18n.Bundle) {

	homeSrv := home.NewHomeService()
	publicHdl := PublicHandler.NewPublicHandler(homeSrv, i18nBundle)

	r.GET("/", publicHdl.Index)

	//find by sku (detail of  a product )
	r.GET("/:category_slug/:product_slug/:sku", publicHdl.ShowProduct)

	r.GET("/:category_slug", publicHdl.ShowProductsByCategory)

	r.GET("/login", publicHdl.ShowLogin)
	r.GET("/verify", publicHdl.ShowVerifyOtp)

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
