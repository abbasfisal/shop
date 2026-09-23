package requests

type StoreAddressRequest struct {
	ReceiverName       string `form:"receiver_name" binding:"required"`
	ReceiverMobile     string `form:"receiver_mobile" binding:"required"`
	ReceiverAddress    string `form:"receiver_address" binding:"required"`
	ReceiverPostalCode string `form:"receiver_postal_code" binding:"required"`
}
