package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/domain/domain_err"
	"shop/interfaces/http/response"
	"shop/pkg/sessions"
)

func (a *AdminHandler) IndexCustomer(c *gin.Context) {
	customers, err := a.customerSrv.Index(c)
	if err.Code == 404 {
		sessions.Set(c, "message", domain_err.RecordNotFound)
	}
	if err.Code == 500 {
		sessions.Set(c, "message", domain_err.InternalServerError)
	}

	fmt.Println("----- customers data : ", customers)
	response.Render(c, http.StatusOK, "admin_index_customer", gin.H{
		"TITLE":     "مدیریت مشتریان",
		"CUSTOMERS": customers,
	})

}
