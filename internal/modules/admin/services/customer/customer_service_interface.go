package customer

import (
	"github.com/gin-gonic/gin"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type CustomerServiceInterface interface {
	Index(c *gin.Context) (responses.Customers, custom_error.CustomError)
}
