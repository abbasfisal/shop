package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"golang.org/x/time/rate"
	"os"
	"shop/application/usecases/attribute"
	attributeValue "shop/application/usecases/attribute_value"
	"shop/application/usecases/auth"
	"shop/application/usecases/banner"
	"shop/application/usecases/brand"
	"shop/application/usecases/category"
	"shop/application/usecases/customer"
	"shop/application/usecases/dashboard"
	order "shop/application/usecases/order"
	"shop/application/usecases/product"
	"shop/bootstrap"
	"shop/infrastructure/database/postgres"
	attributeRepository "shop/infrastructure/repositories/attribute"
	attributeValueRepository "shop/infrastructure/repositories/attribute_value"
	authRepository "shop/infrastructure/repositories/auth"
	bannerRepository "shop/infrastructure/repositories/banner"
	brandRepository "shop/infrastructure/repositories/brand"
	categoryRepository "shop/infrastructure/repositories/category"
	customerRepository "shop/infrastructure/repositories/customer"
	dashboardRepository "shop/infrastructure/repositories/dashboard"
	orderRepository "shop/infrastructure/repositories/order"
	productRepository "shop/infrastructure/repositories/product"
	AdminHandler "shop/interfaces/http/handlers/admin"
	"shop/interfaces/http/middleware"
	"time"
)

func SetAdminRoutes(r *gin.Engine, dep *bootstrap.Dependencies) {

	authRepo := authRepository.NewAuthenticateRepository(postgres.Get())
	authSrv := auth.NewAuthenticateService(authRepo)

	categoryRepo := categoryRepository.NewCategoryRepository(postgres.Get())
	categorySrv := category.NewCategoryService(categoryRepo)

	productRepo := productRepository.NewProductRepository(postgres.Get())
	productSrv := product.NewProductService(productRepo)

	attributeRep := attributeRepository.NewAttributeRepository(postgres.Get())
	attributeSrv := attribute.NewAttributeService(attributeRep)

	attributeValueRepo := attributeValueRepository.NewAttributeRepository(postgres.Get())
	attributeValueSrv := attributeValue.NewAttributeValueService(attributeValueRepo)

	brandRepo := brandRepository.NewBrandRepository(postgres.Get())
	brandSrv := brand.NewBrandService(brandRepo)

	customerRepo := customerRepository.NewCustomerRepository(postgres.Get())
	customerSrv := customer.NewCustomerService(customerRepo)

	orderRepo := orderRepository.NewOrderRepository(postgres.Get())
	orderSrv := order.NewOrderService(orderRepo)

	dashboardSrv := dashboard.NewDashboardService(dashboardRepository.NewDashboardRepository(postgres.Get()))

	bannerSrv := banner.NewBannerService(bannerRepository.NewBannerRepository(postgres.Get()))

	adminHlr := AdminHandler.NewAdminHandler(authSrv, categorySrv, productSrv, attributeSrv, attributeValueSrv, brandSrv, customerSrv, orderSrv, dashboardSrv, bannerSrv, dep)

	// rate limiter
	limiter := middleware.NewRateLimiter(rate.Every(time.Minute), 5)

	guestGrp := r.Group("/")
	guestGrp.Use(middleware.IsGuest, limiter.Middleware())
	{
		guestGrp.GET("/admins/login", adminHlr.ShowLogin)
		guestGrp.POST("/admins/login", adminHlr.PostLogin)
	}

	authGrp := r.Group("/")
	authGrp.Use(middleware.IsAdmin)
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
