package requests

type LoginRequest struct {
	Mobile   string `form:"mobile" binding:"required"`
	Password string `form:"password" binding:"required"`
}
