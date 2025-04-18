package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"golang.org/x/time/rate"
	"os"
	"shop/internal/database/mongodb"
	"shop/internal/database/mysql"
	"shop/internal/middlewares"
	AdminHandler "shop/internal/modules/admin/handlers"
	attributeRepository "shop/internal/modules/admin/repositories/attribute"
	attributeValueRepository "shop/internal/modules/admin/repositories/attribute_value"
	authRepository "shop/internal/modules/admin/repositories/auth"
	bannerRepository "shop/internal/modules/admin/repositories/banner"
	brandRepository "shop/internal/modules/admin/repositories/brand"
	categoryRepository "shop/internal/modules/admin/repositories/category"
	customerRepository "shop/internal/modules/admin/repositories/customer"
	dashboardRepository "shop/internal/modules/admin/repositories/dashboard"
	orderRepository "shop/internal/modules/admin/repositories/order"
	productRepository "shop/internal/modules/admin/repositories/product"
	"shop/internal/modules/admin/services/attribute"
	attributeValue "shop/internal/modules/admin/services/attribute_value"
	"shop/internal/modules/admin/services/auth"
	"shop/internal/modules/admin/services/banner"
	"shop/internal/modules/admin/services/brand"
	"shop/internal/modules/admin/services/category"
	"shop/internal/modules/admin/services/customer"
	"shop/internal/modules/admin/services/dashboard"
	order "shop/internal/modules/admin/services/order"
	"shop/internal/modules/admin/services/product"
	"shop/internal/pkg/bootstrap"
	"time"
)

func SetAdminRoutes(r *gin.Engine, dep *bootstrap.Dependencies) {

	authRepo := authRepository.NewAuthenticateRepository(mysql.Get())
	authSrv := auth.NewAuthenticateService(authRepo)

	categoryRepo := categoryRepository.NewCategoryRepository(mysql.Get())
	categorySrv := category.NewCategoryService(categoryRepo)

	productRepo := productRepository.NewProductRepository(mysql.Get(), mongodb.Get())
	productSrv := product.NewProductService(productRepo)

	attributeRep := attributeRepository.NewAttributeRepository(mysql.Get())
	attributeSrv := attribute.NewAttributeService(attributeRep)

	attributeValueRepo := attributeValueRepository.NewAttributeRepository(mysql.Get())
	attributeValueSrv := attributeValue.NewAttributeValueService(attributeValueRepo)

	brandRepo := brandRepository.NewBrandRepository(mysql.Get())
	brandSrv := brand.NewBrandService(brandRepo)

	customerRepo := customerRepository.NewCustomerRepository(mysql.Get())
	customerSrv := customer.NewCustomerService(customerRepo)

	orderRepo := orderRepository.NewOrderRepository(mysql.Get(), mongodb.Get())
	orderSrv := order.NewOrderService(orderRepo)

	dashboardSrv := dashboard.NewDashboardService(dashboardRepository.NewDashboardRepository(mysql.Get()))

	bannerSrv := banner.NewBannerService(bannerRepository.NewBannerRepository(mysql.Get()))

	adminHlr := AdminHandler.NewAdminHandler(authSrv, categorySrv, productSrv, attributeSrv, attributeValueSrv, brandSrv, customerSrv, orderSrv, dashboardSrv, bannerSrv, dep)

	// rate limiter
	limiter := middlewares.NewRateLimiter(rate.Every(time.Minute), 5)

	guestGrp := r.Group("/")
	guestGrp.Use(middlewares.IsGuest, limiter.Middleware())
	{
		guestGrp.GET("/admins/login", adminHlr.ShowLogin)
		guestGrp.POST("/admins/login", adminHlr.PostLogin)
	}

	authGrp := r.Group("/")
	authGrp.Use(middlewares.IsAdmin)
	{

		//----- asynq monitor panel
		h := asynqmon.New(asynqmon.Options{
			RootPath:     "/admins/monitoring",
			RedisConnOpt: asynq.RedisClientOpt{Addr: fmt.Sprintf(":%s", os.Getenv("REDIS_PORT"))},
		})
		authGrp.Any(h.RootPath()+"/*any", gin.WrapH(h))
		//--------------------------------------------------------

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
		authGrp.POST("/admins/products/:id", adminHlr.UpdateProduct)

		authGrp.GET("/admins/products/:id/add-feature", adminHlr.CreateProductFeature)
		authGrp.POST("/admins/products/:id/add-feature", adminHlr.StoreProductFeature)
		authGrp.GET("/admins/products/:id/show-feature", adminHlr.ShowProductFeature)
		authGrp.GET("/admins/products/:id/delete-feature/:featureID", adminHlr.DeleteProductFeature)
		authGrp.GET("/admins/products/:id/edit-feature/:featureID", adminHlr.EditProductFeature)
		authGrp.POST("/admins/products/:id/update-feature/:featureID", adminHlr.UpdateProductFeature)

		authGrp.GET("/admins/products/:id/show-gallery", adminHlr.ShowProductGallery)
		authGrp.GET("/admins/products/images/:id/delete", adminHlr.DeleteProductImage)
		authGrp.POST("/admins/products/:id/add-images", adminHlr.UploadProductImages)

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

		//order
		authGrp.GET("/admins/orders", adminHlr.IndexOrders)
		authGrp.GET("/admins/orders/:id/details", adminHlr.ShowOrder)
		authGrp.POST("/admins/orders/:id/update-status", adminHlr.EditOrder)

		//banner
		authGrp.GET("/admins/banners/create", adminHlr.CreateBanner)
		authGrp.POST("/admins/banners", adminHlr.StoreBanner)

	}

}
