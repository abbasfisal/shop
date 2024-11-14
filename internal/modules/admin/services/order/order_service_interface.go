package order

import (
	"github.com/gin-gonic/gin"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/pagination"
)

type OrderServiceInterface interface {
	GetOrderPaginate(c *gin.Context) (pagination.Pagination, error)
	GetOrderBy(c *gin.Context, orderID int) (responses.OrderDetail, error)
}
