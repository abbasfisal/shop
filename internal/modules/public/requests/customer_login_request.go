package requests

type CustomerLoginRequest struct {
	Mobile string `form:"mobile" binding:"required"`
}
