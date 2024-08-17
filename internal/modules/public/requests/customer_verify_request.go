package requests

type CustomerVerifyRequest struct {
	N1 string `form:"n1" binding:"required"`
	N2 string `form:"n2" binding:"required"`
	N3 string `form:"n3" binding:"required"`
	N4 string `form:"n4" binding:"required"`
}
