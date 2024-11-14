package responses

import "shop/internal/entities"

type Payment struct {
	CustomerID        uint
	OrderID           uint
	Authority         string
	Description       string
	PaymentURL        string
	StatusCode        int
	Amount            uint
	RefID             string
	Status            int //payment status -> 0-pending,1-paid ,2-failed
	PaymentStatusText string
}

func ToPayment(paymentEntity entities.Payment) Payment {
	return Payment{
		CustomerID:        paymentEntity.CustomerID,
		OrderID:           paymentEntity.OrderID,
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
	case 0:
		return "در حال پرداخت"
	case 1:
		return "پرداخت شده"
	case 2:
		return "لغو شده"
	default:
		return "نامعلوم"
	}

}
