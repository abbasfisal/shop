package handlers

import (
	"context"
	errors2 "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"

	responses "shop/application/dto/admin"
	"shop/application/usecases/banner"
	"shop/application/usecases/fee"
	"shop/application/usecases/home"
	sliders "shop/application/usecases/product_slider"
	"shop/bootstrap"
	"shop/domain/domain_err"
	"shop/domain/entities"
	"shop/interfaces/http/response"
	"shop/pkg/helpers"
	"shop/pkg/util"
)

type PublicHandler struct {
	homeSrv    home.HomeServiceInterface
	bannerSrv  *banner.BannerService
	slidersSrv *sliders.ProductSliderService
	feeSrv     *fee.FeeRateService
	dep        *bootstrap.Dependencies
}

func NewPublicHandler(
	homeSrv home.HomeServiceInterface,
	bannerSrv *banner.BannerService,
	slidersSrv *sliders.ProductSliderService,
	feeSrv *fee.FeeRateService,
	dep *bootstrap.Dependencies,
) PublicHandler {
	return PublicHandler{
		homeSrv:    homeSrv,
		bannerSrv:  bannerSrv,
		slidersSrv: slidersSrv,
		feeSrv:     feeSrv,
		dep:        dep,
	}
}

func (p PublicHandler) Index(c *gin.Context) {
	//20 latest
	//20 random /-
	//20 lowest quantity
	//10 category
	//filter by price
	//filter by quantity
	//pagination in selected category

	products, err := p.homeSrv.GetProducts(context.TODO(), 20)
	if err.Code == 404 {
		c.JSON(200, gin.H{"err": 404})
		return
	}
	if err.Code == 500 {
		c.JSON(200, gin.H{"err": 500})
		return
	}

	categories, cErr := p.homeSrv.GetCategories(context.TODO(), 20)
	if cErr.Code == 404 {
		c.JSON(200, gin.H{"err": 404})
		return
	}
	if cErr.Code == 500 {
		c.JSON(200, gin.H{"err": 500})
		return
	}

	c.JSON(200, gin.H{
		"PRODUCTS":   products,
		"CATEGORIES": categories,
	})
	//	response.Render(c, http.StatusFound, "home", gin.H{
	//			"TITLE":      "home page",
	//			"PRODUCTS":   "",
	//			"CATEGORIES": "",
	//		})
	return
}

func (p PublicHandler) ShowProductsByCategory(c *gin.Context) {

	productPagination, err := p.homeSrv.ListProductByCategorySlug(c, c.Param("category_slug"))
	if err != nil {
		//هر خطایی به جز خطای مرتبط با پیدانکردن رکورد اگر وجود داشت اون خطا رو نشون میدیم
		//در غیر این صورت پیغام رکورد یافت نشد به کاربر نشون داده میشه :)
		if !errors2.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{
				"msg": domain_err.SomethingWrongHappened,
			})
			return
		}
	}

	response.CustomerRender(c, http.StatusFound, "search",
		gin.H{
			"TITLE":          "search",
			"PAGINATION":     productPagination,
			"MEDIA_PATH":     util.GetProductStoragePath(),
			"PrimaryMessage": domain_err.RecordNotFound,
		},
	)
	return

}

// HomePage renders the storefront home: promotion banners (2-up / 4-up) and
// the curated product sliders, each pinned to its homepage slot.
func (p PublicHandler) HomePage(c *gin.Context) {
	bannersTwo := responses.Banners{}
	bannersFour := responses.Banners{}
	if active, err := p.bannerSrv.Active(c); err == nil {
		for _, b := range responses.ToBanners(active).Data {
			if b.Layout == entities.BannerLayoutFour {
				bannersFour.Data = append(bannersFour.Data, b)
			} else {
				bannersTwo.Data = append(bannersTwo.Data, b)
			}
		}
	}

	slidersByPosition := map[string][]responses.ProductSlider{}
	if active, err := p.slidersSrv.ActiveSliders(c); err == nil {
		for position, views := range responses.SlidersByPosition(active) {
			for i := range views {
				views[i].MediaPath = util.GetProductStoragePath()
			}
			slidersByPosition[position] = views
		}
	}

	response.CustomerRender(c, 200, "home", gin.H{
		"TITLE":        "صفحه اصلی فروشگاه",
		"BANNERS_TWO":  bannersTwo,
		"BANNERS_FOUR": bannersFour,
		"SLIDERS":      slidersByPosition,
		"BANNER_PATH":  util.GetBannerStoragePath(),
		"MEDIA_PATH":   util.GetProductStoragePath(),
	})
}

func (p PublicHandler) SingleProduct(c *gin.Context) {

	//get single product with its recommendations product
	product, recommendations, err := p.homeSrv.GetSingleProduct(c, c.Param("product_sku"), c.Param("product_slug"))
	primaryMessage := ""
	if err.Code > 0 {
		if err.Code == 404 {
			primaryMessage = err.DisplayMessage
		} else {
			response.Error500(c)
			return
		}
	}

	response.CustomerRender(c, http.StatusFound, "single_product",
		gin.H{
			"PRODUCT":         product,
			"RECOMMENDATIONS": recommendations,
			"PrimaryMessage":  primaryMessage,
			"MEDIA_PATH":      util.GetProductStoragePath(),
		})
	return
}

func (p PublicHandler) Shipping(c *gin.Context) {
	customer, ok := helpers.GetAuthUser(c)
	if ok {
		if customer.Cart.CartItem.TotalItemCount <= 0 {
			c.Redirect(http.StatusFound, "/checkout/cart")
			return
		}
	}

	// fee quote from the cart items total (the same rule the order stores)
	quote := p.feeSrv.Quote(c, customer.Cart.CartItem.TotalSalePrice)

	response.CustomerRender(c, http.StatusFound, "shipping",
		gin.H{
			"TITLE":           "اطلاعات ارسال",
			"SHIPPING_FEE":    quote.Shipping,
			"PACKAGING_FEE":   quote.Packaging,
			"SHIPPING_FREE":   quote.Free,
			"FREE_THRESHOLD":  quote.Threshold,
			"GRAND_TOTAL":     quote.Grand,
			"SHIPPING_TITLE":  quote.ShipTitle,
			"PACKAGING_TITLE": quote.PackTitle,
		})
	return
}

func (p PublicHandler) ShowOrderList(c *gin.Context) {

	orderPaginations, err := p.homeSrv.ListOrders(c)
	fmt.Println("----------------------orderPagination:", orderPaginations)
	if err != nil || orderPaginations.Rows == nil {
		//هر خطایی به جز خطای مرتبط با پیدانکردن رکورد اگر وجود داشت اون خطا رو نشون میدیم
		//در غیر این صورت پیغام رکورد یافت نشد به کاربر نشون داده میشه :)
		if errors2.Is(err, gorm.ErrRecordNotFound) {
			response.CustomerRender(c, http.StatusNotFound, "profile_orders", gin.H{
				"TITLE":      "لیست سفارشات",
				"PAGINATION": nil,
				"ACTIVE":     "orders",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"msg": domain_err.SomethingWrongHappened,
			})
			return
		}
	}

	response.CustomerRender(c, http.StatusFound, "profile_orders",
		gin.H{
			"TITLE":          "لیست سفارشات",
			"PAGINATION":     orderPaginations,
			"PrimaryMessage": "لیست سفارشات",
			"ACTIVE":         "orders",
		},
	)
	return
}

func (p PublicHandler) ShowOrderDetails(c *gin.Context) {
	q := c.Param("order_number")
	order, err := p.homeSrv.GetOrderBy(c, q)
	if err != nil || order == nil {
		log.Println("---- [public - handlers]-[ShowOrderDetails]----", err)
		c.JSON(http.StatusOK, gin.H{
			"msg": domain_err.SomethingWrongHappened,
		})
		return
	}
	response.CustomerRender(c, http.StatusFound, "customer_order_details", gin.H{
		"TITLE":  "جزییات سفارش",
		"DATA":   order,
		"ACTIVE": "orders",
	})
	return
}

// SliderCatalog is the «مشاهده همه» page of a product slider: every product
// the slider points at (or its whole category scope) with pagination.
func (p PublicHandler) SliderCatalog(c *gin.Context) {
	slider, err := p.slidersSrv.BySlug(c, c.Param("slug"))
	if err != nil || !slider.IsVisibleNow(time.Now()) {
		response.CustomerRender(c, http.StatusNotFound, "404", gin.H{
			"TITLE": "صفحه یافت نشد",
		})
		return
	}

	page, pageErr := p.slidersSrv.Catalog(c, slider)
	view := responses.ToSlider(slider)

	data := gin.H{
		"TITLE":       view.Title,
		"SLIDER":      view,
		"MEDIA_PATH":  util.GetProductStoragePath(),
		"BANNER_PATH": util.GetBannerStoragePath(),
	}
	if pageErr != nil {
		data["PRIMARY_MESSAGE"] = "محصولی در این بخش یافت نشد"
		data["PAGINATION"] = nil
	} else {
		data["PAGINATION"] = page
	}

	response.CustomerRender(c, http.StatusOK, "catalog", data)
	return
}
