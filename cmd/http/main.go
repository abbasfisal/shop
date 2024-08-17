package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"shop/cmd/commands"
	"shop/internal/database/mysql"
	AdminRoutes "shop/internal/modules/admin/routes"
	PublicRoutes "shop/internal/modules/public/routes"
	"shop/internal/pkg/logging"
	"sync"
)

var (
	validate   *validator.Validate
	i18nBundle *i18n.Bundle
	once       sync.Once
	logger     logging.Logger
)

func main() {
	once.Do(initialize)

	commands.Execute()

	r := gin.Default()
	setupSessions(r)
	setupRoutes(r)

	addr := fmt.Sprintf("%s:%s", viper.GetString("App.Host"), viper.GetString("App.Port"))
	if err := r.Run(addr); err != nil {
		logging.GlobalLog.FatalF("[Server start failed]: %v", err)
	}
}

func initialize() {
	loadTranslation()
	initializeLogger()
	loadConfig()
	initializeDatabase()
}

func loadTranslation() {
	i18nBundle = i18n.NewBundle(language.Persian)
	i18nBundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	if _, err := i18nBundle.LoadMessageFile("./internal/translation/active.fa.yaml"); err != nil {
		log.Fatalf("Error loading translation file: %v", err)
	}
}

func initializeLogger() {
	logging.GlobalLog = logging.NewZapLogger()
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config/")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

func initializeDatabase() {
	mysql.Connect()
}

func setupSessions(r *gin.Engine) {
	store := cookie.NewStore([]byte(viper.GetString("App.Key")))
	r.Use(sessions.Sessions("session", store))
}

func setupRoutes(r *gin.Engine) {
	r.LoadHTMLGlob("internal/**/**/**/*.html")
	r.Static("uploads", "./uploads")
	r.Static("assets", "./assets")

	r.StaticFile("/favicon.ico", "./assets/shop/img/seller-logo.png")

	AdminRoutes.SetAdminRoutes(r, i18nBundle)
	PublicRoutes.SetPublic(r, i18nBundle)

	r.GET("/500", func(c *gin.Context) {
		c.HTML(http.StatusInternalServerError, "templates/html/errors/500", nil)
	})
}
