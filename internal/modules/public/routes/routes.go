package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"shop/internal/middlewares"
	PublicHandler "shop/internal/modules/public/handlers"
	CustomerMiddlewares "shop/internal/modules/public/middlewares"
	homeRepository "shop/internal/modules/public/repositories/home"
	"shop/internal/modules/public/repositories/home_mongo"
	"shop/internal/modules/public/services/home"
)

func SetPublic(r *gin.Engine, i18nBundle *i18n.Bundle) {

	//home repository
	homeRep := homeRepository.NewHomeRepository()

	MongoHomeRepo := home_mongo.NewMongoRepository()

	//home service
	homeSrv := home.NewHomeService(homeRep, MongoHomeRepo)

	//--- [Global Middleware]
	r.Use(CustomerMiddlewares.LoadMenu(&homeSrv)) //load menu by LoadMenu middleware
	r.Use(CustomerMiddlewares.CheckUserAuth())    //set `auth` key in context if user was existed in database
	//----

	publicHdl := PublicHandler.NewPublicHandler(homeSrv, i18nBundle)

	r.GET("/", publicHdl.HomePage)
	r.GET("/product/:product_sku/:product_slug", publicHdl.SingleProduct) //show single product
	r.GET("/search/:category_slug", publicHdl.ShowProductsByCategory)     //show products by category

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
		publicAuthGrp.POST("/add-to-cart", publicHdl.AddToCart)
		publicAuthGrp.GET("/checkout/cart", publicHdl.Cart)
		publicAuthGrp.GET("/cart/:cartID/increment", publicHdl.CartItemIncrement)

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
