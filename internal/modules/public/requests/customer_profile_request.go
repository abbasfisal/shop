package requests

type CustomerProfileRequest struct {
	FirstName string `form:"first_name" binding:"required"`
	LastName  string `form:"last_name" binding:"required"`
}
