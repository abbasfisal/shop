package jobs

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"log"
	"net/http/httptest"
	"os"
	"shop/internal/entities"
	"shop/internal/events"
	"shop/internal/modules/public/repositories/home"
	"shop/internal/pkg/bootstrap"
	"shop/internal/pkg/payment/zarinpal"
	"shop/internal/pkg/util"
	"time"
)

const CancelPendingOrders = "orders:pending:cancel"

func TaskCancelPendingOrders() *asynq.Task {
	return asynq.NewTask(CancelPendingOrders, nil)
}

type CancelJob struct {
	*bootstrap.Dependencies
	*events.EventManager
}

func NewCancelJob(dep *bootstrap.Dependencies, em *events.EventManager) *CancelJob {
	return &CancelJob{
		dep,
		em,
	}
}

func (c *CancelJob) ProcessTask(ctx context.Context, t *asynq.Task) error {

	fmt.Println("cancel_pending_orders started")

	var orders []entities.Order

	tenMinutesAgo := time.Now().Add(-10 * time.Minute)
	err := c.DB.WithContext(ctx).
		Table("orders").
		Where("created_at <= ? AND order_status = ?", tenMinutesAgo, entities.OrderPending).
		Preload("OrderItems").
		Preload("Payment").
		Find(&orders).Error

	if err != nil {
		util.Trace("error find pending orders ")
		return err
	}
	if len(orders) == 0 {
		util.Trace("no pending orders found")
		return errors.New("no pending orders found")
	}

	zarin, err := zarinpal.NewZarinpal(os.Getenv("ZARINPAL_MERCHANTID"), false)
	if err != nil {
		log.Println("[zarinpal err]:", err)
		return err
	}

	h := home.NewHomeRepository(c.Dependencies, c.EventManager)

	mockGinCtx, _ := gin.CreateTestContext(httptest.NewRecorder())

	for _, order := range orders {
		if order.Payment.Authority == "" || order.ID == 0 {
			log.Println("authority was empty:", order)
			continue
		}

		go func(order *entities.Order) {
			authority := order.Payment.Authority
			verified, refID, statusCode, vErr := zarin.PaymentVerification(int(order.Payment.Amount), authority)
			fmt.Println("zarinpal payment verification output ", "| verified:", verified, "| refID:", refID, "| statusCode:", statusCode, "| verify error:", vErr)

			if vErr != nil || !verified || (statusCode != 100 && statusCode != 101) {
				return
			}

			util.Trace("call method OrderPaidSuccessfully() ")

			OrderRes, status, errResult := h.OrderPaidSuccessfully(mockGinCtx, order, refID, verified)

			fmt.Println("OrderPaidSuccessfully output", " | orderID:", OrderRes.ID, " | status:", status, "| errorResult:", errResult)

		}(&order)

	}

	return nil
}
