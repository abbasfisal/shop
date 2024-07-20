package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"shop/internal/database/mysql"
	"shop/internal/middlewares"
	AdminHandler "shop/internal/modules/admin/handlers"
	attributeRepository "shop/internal/modules/admin/repositories/attribute"
	attributeValueRepository "shop/internal/modules/admin/repositories/attribute_value"
	authRepository "shop/internal/modules/admin/repositories/auth"
	categoryRepository "shop/internal/modules/admin/repositories/category"
	productRepository "shop/internal/modules/admin/repositories/product"
	"shop/internal/modules/admin/services/attribute"
	attributeValue "shop/internal/modules/admin/services/attribute_value"
	"shop/internal/modules/admin/services/auth"
	"shop/internal/modules/admin/services/category"
	"shop/internal/modules/admin/services/product"
)

func SetAdminRoutes(r *gin.Engine, i18nBundle *i18n.Bundle) {

	authRepo := authRepository.NewAuthenticateRepository(mysql.Get())
	authSrv := auth.NewAuthenticateService(authRepo)

	categoryRepo := categoryRepository.NewCategoryRepository(mysql.Get())
	categorySrv := category.NewCategoryService(categoryRepo)

	productRepo := productRepository.NewProductRepository(mysql.Get())
	productSrv := product.NewProductService(productRepo)

	attributeRep := attributeRepository.NewAttributeRepository(mysql.Get())
	attributeSrv := attribute.NewAttributeService(attributeRep)

	attributeValueRepo := attributeValueRepository.NewAttributeRepository(mysql.Get())
	attributeValueSrv := attributeValue.NewAttributeValueService(attributeValueRepo)

	adminHlr := AdminHandler.NewAdminHandler(authSrv, categorySrv, productSrv, attributeSrv, attributeValueSrv, i18nBundle)

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
		authGrp.GET("/admins/categories/:id/products", adminHlr.CategoryProducts)

		//attributes
		authGrp.GET("/admins/attributes/create", adminHlr.CreateAttribute)
		authGrp.POST("/admins/attributes", adminHlr.StoreAttribute)
		authGrp.GET("/admins/get-attributes/:catID", adminHlr.GetAttributesByCategoryID)

		//attribute-values
		authGrp.GET("/admins/attribute-values/create", adminHlr.CreateAttributeValues)

		//products
		authGrp.GET("/admins/products", adminHlr.IndexProduct)
		authGrp.GET("/admins/products/create", adminHlr.CreateProduct)
		authGrp.POST("/admins/products", adminHlr.StoreProduct)
		authGrp.GET("/admins/products/:id", adminHlr.ShowProduct)
		authGrp.GET("/admins/products/:id/edit", adminHlr.EditProduct)
		authGrp.POST("/admins/products/:id", adminHlr.UpdateProduct)

		//product-attribute
		//product-inventory

	}

}
