package jobs

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"log"
	"shop/internal/entities"
	"shop/internal/pkg/bootstrap"
	"shop/internal/pkg/util"
	"time"
)

const CancelPendingOrders = "orders:pending:cancel"

type OrderItemStruct struct {
	ID          uint
	OrderID     uint
	InventoryID uint
	ProductID   uint
	Quantity    uint
}
type InventoryUpdate struct {
	InventoryID uint
	ProductID   uint
	Quantity    uint
}

func (OrderItemStruct) TableName() string {
	return "order_items" // Define your custom table name here
}

//

func TaskCancelPendingOrders() *asynq.Task {
	return asynq.NewTask(CancelPendingOrders, nil)
}

type CancelJob struct {
	*bootstrap.Dependencies
}

func NewCancelJob(dep *bootstrap.Dependencies) *CancelJob {
	return &CancelJob{
		dep,
	}
}

func (c *CancelJob) ProcessTask(ctx context.Context, t *asynq.Task) error {

	fmt.Println("cancel_pending_orders started")

	type OrderStruct struct {
		ID          uint
		OrderStatus uint
		OrderItems  []*OrderItemStruct `gorm:"foreignKey:OrderID"`
	}

	var orders []OrderStruct
	tenMinutesAgo := time.Now().Add(-2 * time.Minute)
	err := c.Dependencies.DB.WithContext(ctx).
		Table("orders").
		Select("id", "order_status").
		Where("created_at <= ? AND order_status = ?", tenMinutesAgo, entities.OrderPending).
		Preload("OrderItems").
		Limit(1).
		Find(&orders).Error

	if err != nil {
		fmt.Println("err1")
		log.Fatalln("eerr:", err)
	}
	if len(orders) <= 0 {
		fmt.Println("info: there is no any pending order")
		return nil
	}
	fmt.Println("--- order data -----")
	util.PrettyJson(orders)
	// جمع‌آوری اطلاعات لازم برای آپدیت
	var updateInventory []InventoryUpdate
	for _, order := range orders {
		for _, orderItem := range order.OrderItems {
			updateInventory = append(updateInventory, InventoryUpdate{
				InventoryID: orderItem.InventoryID,
				ProductID:   orderItem.ProductID,
				Quantity:    orderItem.Quantity,
			})
		}
	}

	// خواندن inventories  در یک کوئری

	var inventories []entities.ProductInventory
	err = c.Dependencies.DB.WithContext(ctx).
		Where("id IN (?)", getInventoryIDs(updateInventory)).
		Find(&inventories).Error

	if err != nil {
		fmt.Println("err2")
		log.Fatalln("err2 ", err)
	}

	inventoryMap := make(map[uint]entities.ProductInventory)
	for _, inv := range inventories {
		inventoryMap[inv.ID] = inv
	}

	// ایجاد تراکنش
	tx := c.Dependencies.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// مدیریت قفل ردیس
	var lockKeys []string
	for _, update := range updateInventory {
		lockKey := fmt.Sprintf("lock:inventory:%d", update.InventoryID)
		err := retryWithBackoff(3, 100*time.Millisecond, func() error {
			locked, redisErr := c.Dependencies.RedisClient.SetNX(ctx, lockKey, "locked", 2*time.Second).Result()
			if redisErr != nil {
				fmt.Println("err3")
				return redisErr
			}
			if !locked {
				fmt.Println("err4")
				return fmt.Errorf("inventory %d is locked by another process", update.InventoryID)
			}
			lockKeys = append(lockKeys, lockKey)
			return nil
		})

		if err != nil {
			releaseLocks(ctx, c.Dependencies.RedisClient, lockKeys)
			tx.Rollback()
			fmt.Println("err5")
			log.Fatalln("err5:", err)
		}
	}

	// ساخت و اجرای کوئری آپدیت
	query := "UPDATE product_inventories SET reserved_stock = CASE "
	var ids []uint
	for _, update := range updateInventory {
		if inv, exists := inventoryMap[update.InventoryID]; exists {
			newStock := inv.ReservedStock - update.Quantity
			query += fmt.Sprintf("WHEN id = %d THEN %d ", update.InventoryID, newStock)
			ids = append(ids, update.InventoryID)
		}
	}
	query += "END WHERE id IN (?)"

	// تغییر وضعیت سفارش به لغو شده
	for _, order := range orders {
		err := c.Dependencies.DB.WithContext(ctx).Table("orders").Model(&order).
			Where("order_status=?", entities.OrderPending).
			Update("order_status", entities.OrderCancelled).Error
		if err != nil {
			releaseLocks(ctx, c.Dependencies.RedisClient, lockKeys)
			tx.Rollback()
			fmt.Println("err6a")
			log.Fatalln("err6a:", err)
		}
	}

	//todo:
	//change status order , payment
	//sync mongo product product.SyncMongo(ctx,c.Dependencies.DB,)
	if err := tx.Exec(query, ids).Error; err != nil {
		releaseLocks(ctx, c.Dependencies.RedisClient, lockKeys)
		tx.Rollback()
		fmt.Println("err6")

		log.Fatalln("err6:", err)
	}

	if err := tx.Commit().Error; err != nil {
		releaseLocks(ctx, c.Dependencies.RedisClient, lockKeys)
		fmt.Println("err7")

		log.Fatalln("err7:", err)
	}

	releaseLocks(ctx, c.Dependencies.RedisClient, lockKeys)

	return nil
}

// تابع ریترای با بک‌آف
func retryWithBackoff(retries int, delay time.Duration, fn func() error) error {
	for i := 0; i < retries; i++ {
		if err := fn(); err == nil {
			return nil
		}
		time.Sleep(delay)
		delay *= 2 // افزایش تصاعدی تأخیر
	}
	return fmt.Errorf("operation failed after %d retries", retries)
}

// آزادسازی قفل‌ها در Redis
func releaseLocks(ctx context.Context, redisClient *redis.Client, keys []string) {
	for _, key := range keys {
		redisClient.Del(ctx, key)
	}
}

func getInventoryIDs(updates []InventoryUpdate) []uint {
	ids := make([]uint, len(updates))
	for i, u := range updates {
		ids[i] = u.InventoryID
	}
	return ids
}

func getProductIDs(updates []InventoryUpdate) []uint {
	ids := make([]uint, len(updates))
	for i, u := range updates {
		ids[i] = u.ProductID
	}
	return ids
}

//blueprint with HandleXXXTask
//func HandleCancelPendingOrdersTask(ctx context.Context, t *asynq.Task, dependencies bootstrap.Dependencies) error {
//
//	//mysql.Connect()
//	//db := mysql.Get()
//
//	var orders []entities.Order
//
//	tenMinutesAgo := time.Now().Add(-10 * time.Minute)
//	err := db.WithContext(ctx).Where("created_at <=? AND status=?", tenMinutesAgo, entities.OrderPending).Find(&orders).Error
//	if err != nil {
//		fmt.Println("cancel Pending Orders Err:", err)
//	}
//	fmt.Println("--- orders fetched succ --------")
//
//	for _, order := range orders {
//		fmt.Println("order id :", order.ID)
//	}
//
//	return nil
//}
