package responses

import (
	"shop/internal/entities"
	"time"
)

type Payment struct {
	ID                uint
	CustomerID        uint
	OrderID           uint
	CreatedAt         time.Time
	Authority         string
	Description       string
	PaymentURL        string
	StatusCode        int
	Amount            uint
	RefID             string
	Status            int //payment status -> 0-pending,1-paid ,2-failed
	PaymentStatusText string
}

func ToPayment(paymentEntity *entities.Payment) *Payment {
	return &Payment{
		ID:                paymentEntity.ID,
		CustomerID:        paymentEntity.CustomerID,
		OrderID:           paymentEntity.OrderID,
		CreatedAt:         paymentEntity.CreatedAt,
		Authority:         paymentEntity.Authority,
		Description:       paymentEntity.Description,
		PaymentURL:        paymentEntity.PaymentURL,
		StatusCode:        paymentEntity.StatusCode, //passed from web gate
		Amount:            paymentEntity.Amount,
		RefID:             paymentEntity.RefID,
		Status:            paymentEntity.Status,
		PaymentStatusText: PaymentStatusMapper(paymentEntity.Status),
	}
}
func PaymentStatusMapper(status int) string {
	switch status {
	case entities.PaymentPending:
		return "منتظر پرداخت" //0
	case entities.PaymentSuccess:
		return "پرداخت موفق"
	case entities.PaymentFailed:
		return "پرداخت ناموفق یا لغو شده"
	case entities.PaymentRetry:
		return "در انتظار پرداخت مجدد"
	default:
		return "نامعلوم"
	}

}
