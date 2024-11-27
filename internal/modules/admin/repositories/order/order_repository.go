package order

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"shop/internal/entities"
	"shop/internal/modules/admin/requests"
	"shop/internal/pkg/pagination"
	"strconv"
	"strings"
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

	//برای گرفتن دیتای جدول
	//product_attributes
	//مجبور هستیم که ابتدا دو ستون
	//product_id , inventory_id
	//که در جدول order_items هستند
	// رو بدست بیاریم و بعد نتایج اونها رو درون preload استفاده کنیم

	var productAndInventory []struct {
		ProductID   uint
		InventoryID uint
	}
	if err := oRepo.db.WithContext(c).
		Table("order_items").
		Select("product_id , inventory_id").
		Where("order_id = ?", orderID).
		Scan(&productAndInventory).Error; err != nil {
		return order, customer, err
	}

	//حالا نتایج رو به صورت اسلایس در میاریم که مستقیم بشه درون preload استفاده کرد
	var productIDs, inventoryIDs []uint
	for _, item := range productAndInventory {
		productIDs = append(productIDs, item.ProductID)
		inventoryIDs = append(inventoryIDs, item.InventoryID)
	}

	if err := oRepo.db.WithContext(c).
		Preload("OrderItems.Product.ProductInventoryAttributes",
			"product_inventory_attributes.product_id IN (?) AND product_inventory_attributes.product_inventory_id IN (?)",
			productIDs, inventoryIDs,
		).
		Preload("OrderItems.Product.ProductInventoryAttributes.ProductAttribute",
			"product_attributes.product_id IN (?)",
			productIDs,
		).
		Preload("Payment").
		First(&order, orderID).Error; err != nil {
		return order, customer, err
	}

	//select customer
	if err := oRepo.db.WithContext(c).Preload("Address").First(&customer, order.CustomerID).Error; err != nil {
		return order, customer, err
	}

	return order, customer, nil
}

func (oRepo OrderRepository) UpdateOrderStatusAndNote(c *gin.Context, orderID int, req requests.UpdateOrderStatus) error {
	var order entities.Order
	if err := oRepo.db.WithContext(c).Where("id=?", orderID).First(&order).Error; err != nil {
		return gorm.ErrRecordNotFound
	}

	// عدد منفی یک به این معنی هست که که ادمین نمیخواهد حالت پیشفرض وضعیت سفارش را تغییر دهد
	if req.Status != -1 {
		// بررسی معتبر بودن وضعیت سفارش
		//if req.Status < int(entities.OrderPending) || req.Status > int(entities.OrderCompleted) {
		//	return fmt.Errorf("وضعیت سفارش نامعتبر است")
		//}
		order.OrderStatus = uint(req.Status)
	}
	if req.Note != "" {
		order.Note = strings.TrimSpace(req.Note)
	}

	if err := oRepo.db.WithContext(c).Save(&order).Error; err != nil {
		return err
	}
	return nil
}
