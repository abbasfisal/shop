package response

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"shop/application/dto/admin"
	"shop/pkg/converters"
	"shop/pkg/helpers"
	"shop/pkg/sessions"
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
	// OLD is the flat (first value per key) view of OLDS — used by <select> /
	// <radio> state. It must exist on EVERY page: index on a missing key
	// aborts template execution halfway through the response.
	data["OLD"] = oldFields(c)

	user := helpers.Auth(c)
	if user.ID != 0 {
		data["AUTH"] = responses.ToUserResponse(user)
	}

	return data
}

func Error500(c *gin.Context) {
	c.Redirect(http.StatusFound, "/500")
}

/**
*-----------------------------
|		customer render 🛍 مشتری
*-----------------------------
*/

func CustomerRender(c *gin.Context, code int, name string, data gin.H) {
	data = customerWithGlobalData(c, data)

	format := c.DefaultQuery("format", "html")
	if format == "json" {
		c.JSON(code, data)
		return
	}
	c.HTML(code, name, data)
}

func customerWithGlobalData(c *gin.Context, data gin.H) gin.H {
	data["APP_NAME"] = viper.Get("APP.Name")
	data["ERRORS"] = converters.StringToMap(sessions.Flash(c, "errors"))
	data["OLDS"] = converters.StringToUrlValues(sessions.Flash(c, "olds"))
	data["MESSAGE"] = sessions.Flash(c, "message")
	// OLD is the flat (first value per key) view of OLDS — used by <select> /
	// <radio> state. It must exist on EVERY page: index on a missing key
	// aborts template execution halfway through the response.
	data["OLD"] = oldFields(c)

	menu, _ := c.Get("menu") //We load the menu using the LoadMenu() middleware and ignore the ok variable because if there is any error in LoadMenu(), a 500 error will be returned
	data["MENU"] = menu

	//check `auth` key from context
	customer, ok := helpers.GetAuthUser(c)
	if ok && customer.ID > 0 {
		data["AUTH"] = customer
	}

	return data
}

// oldFields flattens the flashed old input to one value per key.
func oldFields(c *gin.Context) map[string]string {
	out := map[string]string{}
	raw := sessions.GET(c, "olds")
	if raw == "" {
		return out
	}
	var form map[string][]string
	if err := json.Unmarshal([]byte(raw), &form); err != nil {
		return out
	}
	for key, vals := range form {
		if len(vals) > 0 {
			out[key] = vals[0]
		}
	}
	return out
}
