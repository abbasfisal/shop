package jobs

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		Where("created_at <= ? AND order_status = ? AND (updated_at >= ? OR created_at <= ?)",
			tenMinutesAgo, entities.OrderPending, tenMinutesAgo, tenMinutesAgo).
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

		if order.Payment == nil || order.Payment.Authority == "" {
			// این حالت به ندرت اتفاق میفته اما احتمال داره که در بعضی حالت ها ما payment نداشته باشیم
			// که دز این حالت صرفا عملیات ازاد کردن موجودی های رزرو شده رو انجام میدیم
			// و بعد وضعیت سفارش و پرداخت رو به حالت کنسل شده تغییر میدیم
			log.Println("authority was empty:", order)
			order.Payment = &entities.Payment{CustomerID: order.CustomerID,
				Description: "payment was not created",
				Authority:   "not_created_" + uuid.New().String(),
			}

			go func() {
				_, _, err := h.OrderPaidSuccessfully(mockGinCtx, &order, "payment was not created", false)
				fmt.Println("OrderPaidSuccessfully err in cancel_pending_orders :", err)
			}()

		} else {

			go func(order *entities.Order) {
				authority := order.Payment.Authority
				verified, refID, statusCode, vErr := zarin.PaymentVerification(int(order.Payment.Amount), authority)
				fmt.Println("zarinpal payment verification output ", "| verified:", verified, "| refID:", refID, "| statusCode:", statusCode, "| verify error:", vErr)

				// read zarinpal status code doc for more information
				// https://www.zarinpal.com/docs/paymentGateway/errorList.html

				// -51 = Session is not valid, session is not active paid try. پرداخت ناموفق
				// 100 = Success عملیات موفق
				// 101 = Verified تراکنش وریفای شده است.
				if statusCode == -51 || statusCode == 101 || statusCode == 100 {

					OrderRes, status, errResult := h.OrderPaidSuccessfully(mockGinCtx, order, refID, verified)
					fmt.Println("OrderPaidSuccessfully output", " | orderID:", OrderRes.ID, " | status:", status, "| errorResult:", errResult)

				} else if vErr != nil || !verified || (statusCode != 100 && statusCode != 101) {
					log.Println(
						"-- somethings wrong happened please check it ",
						" | verified:", verified, " | refID:", refID, " | statusCode:", statusCode, " | vErr:", verified,
					)
					return
				}

			}(&order)
		}

	}

	return nil
}
