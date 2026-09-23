package repositories

import (
	"github.com/gin-gonic/gin"
	"shop/domain/entities"
	"shop/interfaces/http/requests/admin"
	"shop/pkg/pagination"
)

type OrderRepositoryInterface interface {
	GetOrders(c *gin.Context) (pagination.Pagination, error)
	FindOrderBy(c *gin.Context, orderID int) (*entities.Order, *entities.Customer, error)
	UpdateOrderStatusAndNote(c *gin.Context, orderID int, req *requests.UpdateOrderStatus) (*entities.Order, error)
	CancelPendingOrders(c *gin.Context)
}
