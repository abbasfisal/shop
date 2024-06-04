package routes

import (
	"github.com/gin-gonic/gin"
	"shop/internal/middlewares"
	PublicHandler "shop/internal/modules/public/handlers"
)

func SetPublic(r *gin.Engine) {

	publicHdl := PublicHandler.NewPublicHandler()

	r.GET("/", publicHdl.Index)

	//find by sku (detail of  a product )
	r.GET("/:category_slug/:product_slug/:sku", publicHdl.ShowProduct)

	r.GET("/:category_slug", publicHdl.ShowProductsByCategory)

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
