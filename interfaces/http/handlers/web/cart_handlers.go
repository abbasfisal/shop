package handlers

import (
	errors2 "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/domain/domain_err"
	"shop/interfaces/http/requests/web"
	"shop/interfaces/http/response"
	"shop/pkg/errors"
	"shop/pkg/sessions"
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

	// product_id is now the numeric product id (no more MongoDB ObjectID)
	productID, err := strconv.ParseUint(req.ProductID, 10, 64)
	if err != nil {
		fmt.Println("[error]-[AddToCart]:", err)
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	p.homeSrv.AddToCart(c, uint(productID), req)
	c.Redirect(http.StatusFound, c.Request.Referer())
	return
}

func (p PublicHandler) Cart(c *gin.Context) {

	response.CustomerRender(c, 200, "cart",
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
	if errors2.Is(res, domain_err.QuantityExceedsLimit) {
		errors.Init()
		errors.Add(strconv.Itoa(int(req.ProductID)), "سقف سفارش هر محصول ۳ عدد می باشد")
		sessions.Set(c, "errors", errors.ToString())

	}
	if errors2.Is(res, domain_err.OutOfStock) {
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
