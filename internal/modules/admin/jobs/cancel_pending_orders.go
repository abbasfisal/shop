package jobs

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"shop/internal/entities"
	"shop/internal/pkg/bootstrap"
	"time"
)

const CancelPendingOrders = "orders:pending:cancel"

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

	var orders []entities.Order

	tenMinutesAgo := time.Now().Add(-10 * time.Minute)
	err := c.DB.WithContext(ctx).
		Where("created_at <=? AND order_status=?", tenMinutesAgo, entities.OrderPending).
		Find(&orders).
		Error

	if err != nil {
		fmt.Println("cancel Pending Orders Err:", err)
		return err
	}
	fmt.Println("[inf ✅]  orders fetched successful ")

	for _, order := range orders {
		fmt.Println("order id :", order.ID)
	}

	return nil
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
