package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
	"shop/internal/modules/public/services/home"
	"shop/internal/pkg/html"
)

type PublicHandler struct {
	homeSrv    home.HomeServiceInterface
	i18nBundle *i18n.Bundle
}

func NewPublicHandler(homeSrv home.HomeServiceInterface, i18nBundle *i18n.Bundle) PublicHandler {
	return PublicHandler{
		homeSrv:    homeSrv,
		i18nBundle: i18nBundle,
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

func (p PublicHandler) ShowProduct(c *gin.Context) {

	_ = c.Param("category_slug")
	productSlug := c.Param("product_slug")
	sku := c.Param("sku")

	product, err := p.homeSrv.ShowProductDetail(c, productSlug, sku)
	if err.Code == 404 {
		c.Redirect(http.StatusFound, "/")
		return
	}
	if err.Code == 500 {
		html.Error500(c)
		return
	}

	html.Render(c, http.StatusFound, "product_detail", gin.H{
		"PRODUCT": product,
	})
	return
}

func (p PublicHandler) ShowProductsByCategory(c *gin.Context) {

	products, err := p.homeSrv.ShowProductsByCategorySlug(context.TODO(), c.Param("category_slug"))
	if err.Code == 404 {
		c.JSON(200, gin.H{
			"msg": "not found",
		})
		return
	}
	if err.Code == 500 {
		c.JSON(200, gin.H{
			"msg": "internal server err",
		})
		return
	}

	html.Render(c, http.StatusFound, "products_by_category_slug", gin.H{
		"PRODUCTS": products,
	})
	return

}

func (p PublicHandler) ShowLogin(c *gin.Context) {
	html.Render(c, 200, "customer_login", gin.H{
		"TITLE": "اسم فروشگاه",
		"data":  "data",
	})
}

func (p PublicHandler) ShowVerifyOtp(c *gin.Context) {
	html.Render(c, 200, "customer_verify_phone_number", gin.H{
		"TITLE":        "تایید شماره موبایل",
		"PHONE_NUMBER": "093500000000",
	})
}

func (p PublicHandler) HomePage(c *gin.Context) {
	html.Render(c, 200, "home", gin.H{
		"TITLE": "عنوان فروشگاه",
	})
}

func (p PublicHandler) SingleProduct(c *gin.Context) {
	html.Render(c, 200, "customer_single_product", gin.H{
		"TITLE": "عنوان فروشگاه",
	})
}
