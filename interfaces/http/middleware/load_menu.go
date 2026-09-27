package middleware

import (
	"net/http"
	"shop/application/usecases/home"
	"shop/pkg/logging"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
