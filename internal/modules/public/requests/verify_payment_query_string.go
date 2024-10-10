package requests

type VerifyPaymentQueryString struct {
	TransID string `form:"trans_id" binding:"required"`
	OrderID string `form:"order_id" binding:"required"`
	Amount  string `form:"amount" binding:"required"`
}
