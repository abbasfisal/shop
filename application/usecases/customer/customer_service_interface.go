package customer

import (
	"shop/application/dto/admin"
	"shop/domain/domain_err"

	"github.com/gin-gonic/gin"
)

type CustomerServiceInterface interface {
	Index(c *gin.Context) (*responses.Customers, domain_err.CustomError)
}
