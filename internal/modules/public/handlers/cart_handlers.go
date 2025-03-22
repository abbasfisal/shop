package handlers

import (
	errors2 "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"shop/internal/modules/public/requests"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/errors"
	"shop/internal/pkg/html"
	"shop/internal/pkg/sessions"
	"strconv"
)

func (p PublicHandler) AddToCart(c *gin.Context) {
	var req requests.AddToCartRequest

	_ = c.Request.ParseForm()
	bindErr := c.ShouldBind(&req)
	if bindErr != nil {
		c.JSON(200, gin.H{
			"err": bindErr.Error(),
		})
		return
	}
	//validation objectID
	productObjectID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		fmt.Println("[error]-[AddToCart]:", err)
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	p.homeSrv.AddToCart(c, productObjectID, req)
	c.Redirect(http.StatusFound, c.Request.Referer())
	return
}

func (p PublicHandler) Cart(c *gin.Context) {

	html.CustomerRender(c, 200, "cart",
		gin.H{
			"TITLE": "سبد خرید",
		})
	return
}

func (p PublicHandler) CartItemIncrement(c *gin.Context) {
	_ = c.Request.ParseForm()
	var req requests.IncreaseCartItemQty
	err := c.ShouldBind(&req)
	if err != nil {
		c.Redirect(http.StatusFound, "/checkout/cart")
	}

	res := p.homeSrv.CartItemIncrement(c, &req)
	if errors2.Is(res, custom_error.QuantityExceedsLimit) {
		errors.Init()
		errors.Add(strconv.Itoa(int(req.ProductID)), "سقف سفارش هر محصول ۳ عدد می باشد")
		sessions.Set(c, "errors", errors.ToString())

	}
	if errors2.Is(res, custom_error.OutOfStock) {
		errors.Init()
		errors.Add(strconv.Itoa(int(req.ProductID)), "موجودی محصول کافی نمی باشد")
		sessions.Set(c, "errors", errors.ToString())
	}

	c.Redirect(http.StatusFound, "/checkout/cart")
	return
}

func (p PublicHandler) CartItemDecrement(c *gin.Context) {
	_ = c.Request.ParseForm()
	var req requests.IncreaseCartItemQty
	err := c.ShouldBind(&req)
	if err != nil {
		fmt.Println("-- bind error :", err.Error())
		c.Redirect(http.StatusFound, "/checkout/cart")
	}

	p.homeSrv.CartItemDecrement(c, &req)

	c.Redirect(http.StatusFound, "/checkout/cart")
	return
}

func (p PublicHandler) RemoveCartItem(c *gin.Context) {
	_ = c.Request.ParseForm()
	var req requests.IncreaseCartItemQty
	err := c.ShouldBind(&req)
	if err != nil {
		fmt.Println("-- bind error :", err.Error())
		c.Redirect(http.StatusFound, "/checkout/cart")
	}

	p.homeSrv.RemoveCartItem(c, &req)

	c.Redirect(http.StatusFound, "/checkout/cart")
	return
}
