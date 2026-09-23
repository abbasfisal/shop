package requests

type UpdateOrderStatus struct {
	Status int    `form:"status"`
	Note   string `form:"note"`
}
