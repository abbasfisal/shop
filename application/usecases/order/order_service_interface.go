package order

import (
	"shop/application/dto/admin"
	"shop/interfaces/http/requests/admin"
	"shop/pkg/pagination"

	"github.com/gin-gonic/gin"
)

type OrderServiceInterface interface {
	GetOrderPaginate(c *gin.Context) (pagination.Pagination, error)
	GetOrderBy(c *gin.Context, orderID int) (*responses.OrderDetail, error)
	ChangeOrderStatus(c *gin.Context, orderID int, req *requests.UpdateOrderStatus) error
}
