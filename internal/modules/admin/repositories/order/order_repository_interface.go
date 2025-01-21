package order

import (
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
	"shop/internal/modules/admin/requests"
	"shop/internal/pkg/pagination"
)

type OrderRepositoryInterface interface {
	GetOrders(c *gin.Context) (pagination.Pagination, error)
	FindOrderBy(c *gin.Context, orderID int) (*entities.Order, *entities.Customer, error)
	UpdateOrderStatusAndNote(c *gin.Context, orderID int, req *requests.UpdateOrderStatus) (*entities.Order, error)
	CancelPendingOrders(c *gin.Context)
}
