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
	brandRepository "shop/internal/modules/admin/repositories/brand"
	categoryRepository "shop/internal/modules/admin/repositories/category"
	customerRepository "shop/internal/modules/admin/repositories/customer"
	productRepository "shop/internal/modules/admin/repositories/product"
	"shop/internal/modules/admin/services/attribute"
	attributeValue "shop/internal/modules/admin/services/attribute_value"
	"shop/internal/modules/admin/services/auth"
	"shop/internal/modules/admin/services/brand"
	"shop/internal/modules/admin/services/category"
	"shop/internal/modules/admin/services/customer"
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

	brandRepo := brandRepository.NewBrandRepository(mysql.Get())
	brandSrv := brand.NewBrandService(brandRepo)

	customerRepo := customerRepository.NewCustomerRepository(mysql.Get())
	customerSrv := customer.NewCustomerService(customerRepo)

	adminHlr := AdminHandler.NewAdminHandler(authSrv, categorySrv, productSrv, attributeSrv, attributeValueSrv, brandSrv, customerSrv, i18nBundle)

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
		authGrp.GET("/admins/attributes", adminHlr.IndexAttribute)
		authGrp.GET("/admins/attributes/create", adminHlr.CreateAttribute)
		authGrp.POST("/admins/attributes", adminHlr.StoreAttribute)
		authGrp.GET("/admins/attributes/:id/edit", adminHlr.ShowAttribute)
		authGrp.POST("/admins/attributes/:id", adminHlr.UpdateAttribute)

		authGrp.GET("/admins/get-attributes/:catID", adminHlr.GetAttributesByCategoryID)

		//attribute-values
		authGrp.GET("/admins/attribute-values", adminHlr.IndexAttributeValues)
		authGrp.GET("/admins/attribute-values/create", adminHlr.CreateAttributeValues)
		authGrp.POST("/admins/attribute-values", adminHlr.StoreAttributeValues)

		authGrp.GET("/admins/attribute/values/:id/show", adminHlr.ShowAttributeValues) //show attribute-values of an attribute
		authGrp.GET("/admins/attribute-values/:id/edit", adminHlr.EditAttributeValues)
		authGrp.POST("/admins/attribute-values/:id/edit", adminHlr.UpdateAttributeValues)
		///admins/attributes-values/{{.ID}}/edit

		//products
		authGrp.GET("/admins/products", adminHlr.IndexProduct)
		authGrp.GET("/admins/products/create", adminHlr.CreateProduct)
		authGrp.POST("/admins/products", adminHlr.StoreProduct)
		authGrp.GET("/admins/products/:id", adminHlr.ShowProduct)
		authGrp.GET("/admins/products/:id/edit", adminHlr.EditProduct)
		authGrp.GET("/admins/products/:id/show-gallery", adminHlr.ShowProductGallery)
		authGrp.GET("/admins/products/images/:id/delete", adminHlr.DeleteProductImage)
		authGrp.POST("/admins/products/:id/add-images", adminHlr.UploadProductImages)
		authGrp.POST("/admins/products/:id", adminHlr.UpdateProduct)

		//product-attribute
		authGrp.GET("/admins/products/:id/add-attributes", adminHlr.ProductsAddAttributes)
		authGrp.POST("/admins/products/:id/add-attributes", adminHlr.StoreProductsAddAttributes)
		//product-inventory
		authGrp.GET("/admins/products/:id/add-inventory", adminHlr.ShowProductInventory)
		authGrp.POST("/admins/products/:id/add-inventory", adminHlr.StoreProductInventory)
		authGrp.GET("/admins/product-inventory-attributes/:id/delete", adminHlr.DeleteProductInventoryAttribute)
		authGrp.GET("/admins/inventories/:id/delete", adminHlr.DeleteInventory)
		authGrp.POST("/admins/inventories/:id/append-attributes", adminHlr.AppendAttribute)
		authGrp.POST("/admins/inventories/:id/update-quantity", adminHlr.UpdateQuantity)
		//brand
		authGrp.GET("/admins/brands", adminHlr.IndexBrand)
		authGrp.GET("/admins/brands/create", adminHlr.ShowCreateBrand)
		authGrp.POST("/admins/brands/create", adminHlr.StoreBrand)
		authGrp.GET("/admins/brands/:id", adminHlr.ShowBrand)
		authGrp.GET("/admins/brands/:id/edit", adminHlr.EditBrand)
		authGrp.POST("/admins/brands/:id/edit", adminHlr.UpdateBrand)

		//customer
		authGrp.GET("/admins/customers", adminHlr.IndexCustomer)

	}

}
