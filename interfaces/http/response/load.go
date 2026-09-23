package response

import "github.com/gin-gonic/gin"

func LoadHtml(r *gin.Engine) {
	//templates/{admin,site,layouts,errors}/*.html
	r.LoadHTMLGlob("templates/*/*.html")
}
