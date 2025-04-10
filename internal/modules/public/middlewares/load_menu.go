package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"shop/internal/modules/public/services/home"
	"shop/internal/pkg/logging"
)

func LoadMenu(homeSrv home.HomeServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		menu, err := homeSrv.GetMenu(c)
		if err != nil {

			logging.Log.
				WithFields(logrus.Fields{"function": "LoadMenu"}).
				WithError(err).Fatal("load menu failed")

			c.HTML(http.StatusInternalServerError, "500", gin.H{})
			c.Abort()

		}
		c.Set("menu", menu)
		c.Next()
	}
}
