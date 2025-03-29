package handlers

import (
	errors2 "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"shop/internal/modules/admin/requests"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/html"
	"shop/internal/pkg/sessions"
	"strconv"
)

func (a *AdminHandler) IndexOrders(c *gin.Context) {
	orderPaginate, err := a.orderSrv.GetOrderPaginate(c)

	if err != nil {
		if errors2.Is(err, gorm.ErrRecordNotFound) {
			html.Render(c, http.StatusFound, "admin_index_order", gin.H{
				"TITLE":           "لیست سفارشات",
				"PRIMARY_MESSAGE": "سفارشی موجود نیست",
				"PAGINATION":      nil,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	html.Render(c, http.StatusFound, "admin_index_order", gin.H{
		"TITLE":      "لیست سفارشات",
		"PAGINATION": orderPaginate,
	})
	return
}

func (a *AdminHandler) ShowOrder(c *gin.Context) {

	//get and convert order id from url
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", "ID سفارش صحیح نمی باشد.")
		c.Redirect(http.StatusFound, "/admins/orders")
		return
	}

	orderRes, err := a.orderSrv.GetOrderBy(c, orderID)
	if err != nil {
		sessions.Set(c, "message", err.Error())
		c.Redirect(http.StatusFound, "/admins/orders")
		return
	}

	html.Render(c, http.StatusOK, "admin_show_order", gin.H{
		"TITLE":    "جزییات سفارش",
		"Customer": orderRes.Customer,
		"Data":     orderRes.Order,
	})
	return
}

func (a *AdminHandler) EditOrder(c *gin.Context) {
	//get order id from url
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admins/orders")
		return
	}

	//bind request
	var req requests.UpdateOrderStatus
	bindErr := c.ShouldBind(&req)
	if bindErr != nil {
		c.JSON(200, gin.H{
			"message": "bind error ",
			"err":     err,
		})
		return
	}

	url := fmt.Sprintf("/admins/orders/%d/details", orderID)

	//اگه ادمین دیتای خالی ارسال کنه برای مدل سفارش اتفاقی نباید بیفته
	if req.Note == "" && req.Status == -1 {
		c.Redirect(http.StatusFound, url)
		return
	}

	//update order
	if err := a.orderSrv.ChangeOrderStatus(c, orderID, &req); err != nil {
		sessions.Set(c, "message", custom_error.UpdateOrderFaileds)
		c.Redirect(http.StatusFound, url)
		return
	}

	sessions.Set(c, "message", custom_error.SuccessfullyUpdated)
	c.Redirect(http.StatusFound, url)
	return
}
