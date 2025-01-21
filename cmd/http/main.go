package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
	"shop/cmd/commands"
	"shop/internal/database/mongodb"
	"shop/internal/database/mysql"
	"shop/internal/events"
	AdminRoutes "shop/internal/modules/admin/routes"
	PublicRoutes "shop/internal/modules/public/routes"
	"shop/internal/pkg/bootstrap"
	"shop/internal/pkg/cache"
	"shop/internal/pkg/logging"
	"shop/internal/pkg/util"
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

	r.SetFuncMap(template.FuncMap{
		"stringToUint": util.StringToUint,
	})

	setupSessions(r)
	setupRoutes(r, dependencies.I18nBundle, dependencies.AsynqClient, dependencies)

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

func setupRoutes(r *gin.Engine, i18nBundle *i18n.Bundle, asynqClient *asynq.Client, dep *bootstrap.Dependencies) {

	r.LoadHTMLGlob("internal/**/**/**/*.html")
	r.Static("uploads", "./uploads")
	r.Static("assets", "./assets")
	r.StaticFile("/favicon.ico", "./assets/shop/img/seller-logo.png")

	//eventManager dependencies : note: we have to prevent cycle import ,
	//									therefor we create another dependency struct
	eventManagerDep := events.EventManagerDep{
		AsynqClient: asynqClient,
		DB:          mysql.Get(),
		RedisClient: cache.NewRedisClient(),
		MongoClient: mongodb.Get(),
	}
	em := events.NewEventManager(&eventManagerDep)

	AdminRoutes.SetAdminRoutes(r, i18nBundle, asynqClient)
	PublicRoutes.SetPublic(r, dep, em)

	r.GET("/500", func(c *gin.Context) {
		c.HTML(http.StatusInternalServerError, "templates/html/errors/500", nil)
	})
}
