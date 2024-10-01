package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"shop/internal/pkg/helpers"
)

// CheckUserAuth چک میکنه که ایا کاربر با این سشن ای دیش وجود داره یا نه اگر وجود داشت درون کانتکس یک کلید ست میکنه
func CheckUserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("[middleware] :CheckUserAuth")

		customer := helpers.CustomerAuth(c) //check session id and check user exist with passed session or not
		if customer.ID > 0 {
			c.Set("auth", customer) //yes user was existed , so we set a key `auth` in context
		}

		c.Next()
	}
}
