package order

import (
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
	"shop/internal/pkg/pagination"
)

type OrderRepositoryInterface interface {
	GetOrders(c *gin.Context) (pagination.Pagination, error)
	FindOrderBy(c *gin.Context, orderID int) (entities.Order, entities.Customer, error)
}
