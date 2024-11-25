package order

import (
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
	"shop/internal/modules/admin/repositories/order"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/pagination"
)

type OrderService struct {
	repo order.OrderRepositoryInterface
}

func NewOrderService(repo order.OrderRepositoryInterface) OrderService {
	return OrderService{repo: repo}
}

func (o OrderService) GetOrderPaginate(c *gin.Context) (pagination.Pagination, error) {

	orderList, err := o.repo.GetOrders(c)
	if err != nil {
		return pagination.Pagination{}, err
	}

	orderList.Rows = responses.ToOrders(orderList.Rows.([]entities.Order))
	return orderList, nil
}

func (o OrderService) GetOrderBy(c *gin.Context, orderID int) (responses.OrderDetail, error) {
	orderEntity, customerEntity, err := o.repo.FindOrderBy(c, orderID)
	if err != nil {
		return responses.OrderDetail{}, err
	}

	return responses.ToOrderDetail(orderEntity, customerEntity), nil
}