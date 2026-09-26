package routes

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	bannerUseCase "shop/application/usecases/banner"
	feeUseCase "shop/application/usecases/fee"
	"shop/application/usecases/home"
	sliderUseCase "shop/application/usecases/product_slider"
	"shop/bootstrap"
	"shop/infrastructure/database/postgres"
	"shop/infrastructure/events"
	bannerRepository "shop/infrastructure/repositories/banner"
	feeRepository "shop/infrastructure/repositories/fee"
	mysqlRepo "shop/infrastructure/repositories/home"
	sliderRepository "shop/infrastructure/repositories/product_slider"
	PublicHandler "shop/interfaces/http/handlers/web"
	"shop/interfaces/http/middleware"
	"time"
)

func SetPublic(r *gin.Engine, dep *bootstrap.Dependencies, eventManager *events.EventManager) {
	// note: we need to access to the eventManager in everywhere like repo , service
	repo := mysqlRepo.NewHomeRepository(dep, eventManager)

	//home service
	homeSrv := home.NewHomeService(dep, repo, eventManager)

	//--- [Global Middleware]
	publicLimiter := middleware.NewRateLimiter(rate.Every(time.Minute), 60)
	limiter := middleware.NewRateLimiter(rate.Every(time.Minute), 5) // use for specific routes

	r.Use(middleware.LoadMenu(homeSrv)) //load menu by LoadMenu middleware
	r.Use(middleware.CheckUserAuth())   //set `auth` key in context if user was existed in database
	r.Use(publicLimiter.Middleware())
	//----

	// promotion banners + homepage product sliders (storefront feed)
	bannerSrv := bannerUseCase.NewBannerService(bannerRepository.NewBannerRepository(postgres.Get()))
	sliderSrv := sliderUseCase.NewProductSliderService(sliderRepository.NewProductSliderRepository(postgres.Get()))
	feeSrv := feeUseCase.NewFeeRateService(feeRepository.NewFeeRateRepository(postgres.Get()))

	publicHdl := PublicHandler.NewPublicHandler(homeSrv, bannerSrv, sliderSrv, feeSrv, dep)

	r.GET("/", publicHdl.HomePage)
	r.GET("/product/:product_sku/:product_slug", publicHdl.SingleProduct) //show single product
	r.GET("/search/:category_slug", publicHdl.ShowProductsByCategory)     //show products by category
	r.GET("/sliders/:slug", publicHdl.SliderCatalog)                      //catalog page of one product slider
	r.GET("/checkout/payment/verify", publicHdl.VerifyPayment)            //payment callback url

	r.GET("/tsearch", publicHdl.SearchProductByTypesence) //search product with typesence
	r.GET("/tsearch/show", publicHdl.ShowTypeSenceForm)   //show typesence html form

	publicAuthGrp := r.RouterGroup
	customerRoute := r.RouterGroup
	customerRoute.Use(limiter.Middleware(), middleware.CheckCustomerSessionID())
	{
		customerRoute.GET("/login", publicHdl.ShowLogin)
		customerRoute.POST("/login", publicHdl.PostLogin)
		customerRoute.GET("/verify", publicHdl.ShowVerifyOtp)
		customerRoute.POST("/verify", publicHdl.PostVerifyOtp)
		customerRoute.GET("/resend-otp", publicHdl.ResendOtp)
	}

	publicAuthGrp.Use(middleware.CustomerMustLogin())
	{
		publicAuthGrp.GET("/logout", publicHdl.LogOut)
		publicAuthGrp.GET("/profile", publicHdl.ShowProfile)
		publicAuthGrp.GET("/profile/edit", publicHdl.EditProfile)
		publicAuthGrp.POST("/profile/edit", publicHdl.UpdateProfile)

		publicAuthGrp.POST("/add-to-cart", publicHdl.AddToCart)            //insert
		publicAuthGrp.GET("/checkout/cart", publicHdl.Cart)                //get-all
		publicAuthGrp.POST("/cart/increment", publicHdl.CartItemIncrement) //+
		publicAuthGrp.POST("/cart/decrement", publicHdl.CartItemDecrement) //-
		publicAuthGrp.POST("/cart/remove", publicHdl.RemoveCartItem)       //delete

		publicAuthGrp.GET("/checkout/shipping", publicHdl.Shipping)    //shipping
		publicAuthGrp.POST("/addresses/store", publicHdl.StoreAddress) //store address

		publicAuthGrp.POST("/checkout/payment", publicHdl.Payment) //payment

		//-- orders
		publicAuthGrp.GET("/orders", publicHdl.ShowOrderList)
		publicAuthGrp.GET("/orders/detail/:order_number", publicHdl.ShowOrderDetails)

	}

	guestGrp := r.Group("/")
	guestGrp.Use(middleware.IsGuest)
	{
	}
}
