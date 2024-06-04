package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"shop/cmd/commands"
	"shop/internal/database/mysql"
	AdminRoutes "shop/internal/modules/admin/routes"
	PublicRoutes "shop/internal/modules/public/routes"
)

func main() {

	//load config
	configInit()

	//load mysql connection
	mysql.Connect()

	commands.Execute()

	r := gin.Default()

	store := cookie.NewStore([]byte(viper.GetString("App.Key")))
	r.Use(sessions.Sessions("session", store))

	r.LoadHTMLGlob("internal/**/**/**/*.html")
	r.Static("uploads", "./uploads")
	r.Static("assets", "./assets")
	//Admin routes
	AdminRoutes.SetAdminRoutes(r)

	//public routes
	PublicRoutes.SetPublic(r)

	r.GET("/500", func(c *gin.Context) {
		c.HTML(http.StatusInternalServerError, "templates/html/errors/500", nil)
		return
	})

	if err := r.Run(fmt.Sprintf("%s:%s", viper.GetString("App.Host"), viper.GetString("App.Port"))); err != nil {
		log.Fatal("[Server start failed ] : ", err)
	}

}

func configInit() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("error reading config file ", err)
	}
}
