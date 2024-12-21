package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"shop/cmd/commands"
	"shop/internal/database/mongodb"
	"shop/internal/database/mysql"
	AdminRoutes "shop/internal/modules/admin/routes"
	PublicRoutes "shop/internal/modules/public/routes"
	"shop/internal/pkg/cache"
	"shop/internal/pkg/logging"
	"shop/internal/pkg/sms"
	"sync"
)

var (
	i18nBundle  *i18n.Bundle
	asynqClient *asynq.Client
	once        sync.Once
)

func main() {
	once.Do(initialize)

	defer mysql.Close()

	cache.InitRedisClient()

	mongodb.Connect()
	defer mongodb.Disconnect()

	commands.Execute()

	r := gin.Default()
	setupSessions(r)
	setupRoutes(r, asynqClient)

	addr := fmt.Sprintf("%s:%s", viper.GetString("App.Host"), viper.GetString("App.Port"))
	log.Printf("[start server ]: %s", "http://"+addr)
	if err := r.Run(addr); err != nil {
		logging.GlobalLog.FatalF("[Server start failed]: %v", err)
	}
}

func initialize() {
	loadConfig()
	loadTranslation()
	initializeLogger()
	initializeDatabase()
	InitializeSmsService()
	initAsynqClient()
}

func initAsynqClient() {
	opt := asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_DB"), os.Getenv("REDIS_PORT"))}
	asynqClient = asynq.NewClient(opt)
}

func InitializeSmsService() {
	kaveNegar := sms.NewKaveNegar(os.Getenv("KAVENEGAR_SECRETKEY"))
	sms.GetSMSManager().SetService(kaveNegar)
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
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("--[main.go]-[-- load env failed --]- Warning: No .env file found or unable to load: %v", err)
	}
	viper.AutomaticEnv() //read local environment automatically

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

func setupRoutes(r *gin.Engine, asynqClient *asynq.Client) {
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
