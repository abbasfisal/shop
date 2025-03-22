package handlers

import (
	"context"
	errors2 "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"shop/internal/modules/public/services/home"
	"shop/internal/pkg/bootstrap"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/helpers"
	"shop/internal/pkg/html"
)

type PublicHandler struct {
	homeSrv home.HomeServiceInterface
	dep     *bootstrap.Dependencies
}

func NewPublicHandler(homeSrv home.HomeServiceInterface, dep *bootstrap.Dependencies) PublicHandler {
	return PublicHandler{
		homeSrv: homeSrv,
		dep:     dep,
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
	//	html.Render(c, http.StatusFound, "home", gin.H{
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
				"msg": custom_error.SomethingWrongHappened,
			})
			return
		}
	}

	html.CustomerRender(c, http.StatusFound, "search",
		gin.H{
			"TITLE":          "search",
			"PAGINATION":     productPagination,
			"PrimaryMessage": custom_error.RecordNotFound,
		},
	)
	return

}

func (p PublicHandler) HomePage(c *gin.Context) {
	//menu, err := p.homeSrv.GetMenu(c)
	//if err != nil {
	//	html.Error500(c)
	//	return
	//}
	//row-header
	//row-newest
	//row-random
	//row-banners
	//row-by-category
	//

	html.CustomerRender(c, 200, "home", gin.H{
		"TITLE": "صفحه اصلی فروشگاه",
	})
}

func (p PublicHandler) SingleProduct(c *gin.Context) {

	product, err := p.homeSrv.GetSingleProduct(c, c.Param("product_sku"), c.Param("product_slug"))
	primaryMessage := ""
	if err.Code > 0 {
		if err.Code == 404 {
			primaryMessage = err.DisplayMessage
		} else {
			html.Error500(c)
			return
		}
	}

	html.CustomerRender(c, http.StatusFound, "single_product",
		gin.H{
			"PRODUCT":        product,
			"PrimaryMessage": primaryMessage,
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

	html.CustomerRender(c, http.StatusFound, "shipping",
		gin.H{
			"TITLE": "اطلاعات ارسال",
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
			html.CustomerRender(c, http.StatusNotFound, "profile_orders", gin.H{
				"TITLE":      "لیست سفارشات",
				"PAGINATION": nil,
				"ACTIVE":     "orders",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"msg": custom_error.SomethingWrongHappened,
			})
			return
		}
	}

	html.CustomerRender(c, http.StatusFound, "profile_orders",
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
			"msg": custom_error.SomethingWrongHappened,
		})
		return
	}
	html.CustomerRender(c, http.StatusFound, "customer_order_details", gin.H{
		"TITLE":  "جزییات سفارش",
		"DATA":   order,
		"ACTIVE": "orders",
	})
	return
}
