package html

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/converters"
	"shop/internal/pkg/helpers"
	"shop/internal/pkg/sessions"
)

func Render(c *gin.Context, code int, name string, data gin.H) {
	data = WithGlobalData(c, data)

	format := c.DefaultQuery("format", "html")
	if format == "json" {
		c.JSON(code, data)
		return
	}
	c.HTML(code, name, data)
}

func WithGlobalData(c *gin.Context, data gin.H) gin.H {
	data["APP_NAME"] = viper.Get("APP.Name")
	data["ERRORS"] = converters.StringToMap(sessions.Flash(c, "errors"))
	data["OLDS"] = converters.StringToUrlValues(sessions.Flash(c, "olds"))
	data["MESSAGE"] = sessions.Flash(c, "message")

	user := helpers.Auth(c)
	if user.ID != 0 {
		data["AUTH"] = responses.ToUserResponse(user)
	}

	return data
}

func Error500(c *gin.Context) {
	c.Redirect(http.StatusFound, "/500")
}
