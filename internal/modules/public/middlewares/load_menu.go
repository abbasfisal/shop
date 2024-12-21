package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/internal/modules/public/services/home"
	"shop/internal/pkg/custom_error"
)

func LoadMenu(homeSrv home.HomeServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		menu, err := homeSrv.GetMenu(c)
		if err != nil {
			fmt.Println("~~~~~~~~ loadMenu Middleware err:", err, "\n~~~~~~~~~~~~~~~~~")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": custom_error.SomethingWrongHappened,
			})
			c.Abort()
			return
		}
		c.Set("menu", menu)
		c.Next()
	}
}
