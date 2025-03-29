package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/html"
	"shop/internal/pkg/sessions"
)

func (a *AdminHandler) IndexCustomer(c *gin.Context) {
	customers, err := a.customerSrv.Index(c)
	if err.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
	}
	if err.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
	}

	fmt.Println("----- customers data : ", customers)
	html.Render(c, http.StatusFound, "admin_index_customer", gin.H{
		"TITLE":     "مدیریت مشتریان",
		"CUSTOMERS": customers,
	})

}
