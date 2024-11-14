package order

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"shop/internal/entities"
	"shop/internal/pkg/pagination"
	"strconv"
)

type OrderRepository struct {
	db          *gorm.DB
	mongoClient *mongo.Client
}

func NewOrderRepository(db *gorm.DB, mongoClient *mongo.Client) OrderRepository {
	return OrderRepository{
		db:          db,
		mongoClient: mongoClient,
	}
}

func (oRepo OrderRepository) GetOrders(c *gin.Context) (pagination.Pagination, error) {

	// Convert query parameters from string to int
	limitStr := c.Query("limit")
	pageStr := c.Query("page")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 { // Default to 10 if invalid
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 { // Default to 1 if invalid
		page = 1
	}
	var pg = pagination.Pagination{
		Limit: limit,
		Page:  page,
	}

	var orders []entities.Order
	//condition := fmt.Sprintf("customer_id=%d", customer.ID)
	condition := ""

	paginateQuery, exist := pagination.Paginate(c, condition, &orders, &pg, oRepo.db)

	if !exist {
		return pg, gorm.ErrRecordNotFound
	}

	//if pErr := paginateQuery(h.db).Preload("OrderItems").Where("customer_id=?", customer.ID).Find(&orders).Error; pErr != nil {
	if pErr := paginateQuery(oRepo.db).Find(&orders).Error; pErr != nil {
		return pg, pErr
	}

	pg.Rows = orders
	return pg, nil
}

func (oRepo OrderRepository) FindOrderBy(c *gin.Context, orderID int) (entities.Order, entities.Customer, error) {
	var order entities.Order
	var customer entities.Customer

	if err := oRepo.db.WithContext(c).Preload("OrderItems").Preload("Payment").First(&order, orderID).Error; err != nil {
		return order, customer, err
	}

	if err := oRepo.db.WithContext(c).First(&customer, order.CustomerID).Error; err != nil {
		return order, customer, err
	}

	return order, customer, nil
}
