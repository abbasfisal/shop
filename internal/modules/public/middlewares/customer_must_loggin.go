package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/internal/pkg/helpers"
)

// CustomerMustLogin  this middleware just check is customer LoggedIn if not redirect
func CustomerMustLogin() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("[Middleware] : CustomerMustLogin")

		customer, ok := helpers.GetAuthUser(c)
		if !ok || customer.ID <= 0 {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
