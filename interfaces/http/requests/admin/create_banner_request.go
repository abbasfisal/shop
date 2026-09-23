package requests

type CreateBannerRequest struct {
	Type        uint   `form:"type" binding:"required"`
	Link        string `form:"link" binding:"required"`
	Priority    uint   `form:"priority" binding:"required"`
	Status      string `form:"status"`
	BannerImage string
}
