package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"shop/cmd/commands"
	"shop/internal/database/mongodb"
	"shop/internal/database/mysql"
	AdminRoutes "shop/internal/modules/admin/routes"
	PublicRoutes "shop/internal/modules/public/routes"
	"shop/internal/pkg/bootstrap"
	"shop/internal/pkg/logging"
)

func main() {

	dependencies, err := bootstrap.Initialize()
	if err != nil {
		log.Fatal("[x] error initializing project :", err)
	}

	defer mysql.Close()
	defer mongodb.Disconnect()
	defer dependencies.AsynqClient.Close()

	commands.Execute()

	r := gin.Default()
	setupSessions(r)
	setupRoutes(r, dependencies.I18nBundle, dependencies.AsynqClient)

	addr := fmt.Sprintf("%s:%s", viper.GetString("App.Host"), viper.GetString("App.Port"))
	log.Printf("[start server ]: %s", "http://"+addr)
	if err := r.Run(addr); err != nil {
		logging.GlobalLog.FatalF("[Server start failed]: %v", err)
	}
}

func setupSessions(r *gin.Engine) {
	store := cookie.NewStore([]byte(viper.GetString("App.Key")))
	r.Use(sessions.Sessions("session", store))
}

func setupRoutes(r *gin.Engine, i18nBundle *i18n.Bundle, asynqClient *asynq.Client) {

	r.LoadHTMLGlob("internal/**/**/**/*.html")
	r.Static("uploads", "./uploads")
	r.Static("assets", "./assets")
	r.StaticFile("/favicon.ico", "./assets/shop/img/seller-logo.png")

	AdminRoutes.SetAdminRoutes(r, i18nBundle, asynqClient)
	PublicRoutes.SetPublic(r, i18nBundle, asynqClient)

	r.GET("/500", func(c *gin.Context) {
		c.HTML(http.StatusInternalServerError, "templates/html/errors/500", nil)
	})
}
