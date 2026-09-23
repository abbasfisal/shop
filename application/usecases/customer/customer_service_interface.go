package customer

import (
	"github.com/gin-gonic/gin"
	"shop/application/dto/admin"
	"shop/domain/domain_err"
)

type CustomerServiceInterface interface {
	Index(c *gin.Context) (*responses.Customers, domain_err.CustomError)
}
