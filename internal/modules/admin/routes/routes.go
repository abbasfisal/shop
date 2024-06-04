package routes

import (
	"github.com/gin-gonic/gin"
	"shop/internal/middlewares"
	AdminHandler "shop/internal/modules/admin/handlers"
)

func SetAdminRoutes(r *gin.Engine) {

	adminHlr := AdminHandler.NewAdminHandler()

	guestGrp := r.Group("/")
	guestGrp.Use(middlewares.IsGuest)
	{
		guestGrp.GET("/admins/login", adminHlr.ShowLogin)
		guestGrp.POST("/admins/login", adminHlr.PostLogin)
		//register routes
	}

	authGrp := r.Group("/")
	authGrp.Use(middlewares.IsAdmin)
	{
		authGrp.GET("/starter", func(c *gin.Context) {
			c.HTML(200, "starter", nil)
			return
		})

		authGrp.GET("/admins/home", adminHlr.ShowHome)

		//categories
		authGrp.GET("/admins/categories", adminHlr.IndexCategory)
		authGrp.GET("/admins/categories/create", adminHlr.CreateCategory)
		authGrp.POST("/admins/categories", adminHlr.StoreCategory)
		authGrp.GET("/admins/categories/:id", adminHlr.ShowCategory)
		authGrp.GET("/admins/categories/:id/edit", adminHlr.EditCategory)
		authGrp.POST("/admins/categories/:id", adminHlr.UpdateCategory)

		//products
		authGrp.GET("/admins/products", adminHlr.IndexProduct)
		authGrp.GET("/admins/products/create", adminHlr.CreateProduct)
		authGrp.POST("/admins/products", adminHlr.StoreProduct)
		authGrp.GET("/admins/products/:id", adminHlr.ShowProduct)
		authGrp.GET("/admins/products/:id/edit", adminHlr.EditProduct)
		authGrp.POST("/admins/products/:id", adminHlr.UpdateProduct)

	}

}
