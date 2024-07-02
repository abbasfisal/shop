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

var validate *validator.Validate
var i18nBundle *i18n.Bundle
var Once sync.Once
var logger logging.Logger

func main() {

	Once.Do(func() {
		//load translation
		loadTranslation()

		//load Logger
		logger = logging.NewZapLogger()

		//load config
		configInit()

		//load mysql connection
		mysql.Connect()
	})
	logger.InfoF("this is just for info")

	commands.Execute()

	r := gin.Default()

	store := cookie.NewStore([]byte(viper.GetString("App.Key")))
	r.Use(sessions.Sessions("session", store))

	r.LoadHTMLGlob("internal/**/**/**/*.html")
	r.Static("uploads", "./uploads")
	r.Static("assets", "./assets")
	//Admin routes
	AdminRoutes.SetAdminRoutes(r, i18nBundle)

	//public routes
	PublicRoutes.SetPublic(r, i18nBundle)

	r.GET("/500", func(c *gin.Context) {
		c.HTML(http.StatusInternalServerError, "templates/html/errors/500", nil)
		return
	})

	if err := r.Run(fmt.Sprintf("%s:%s", viper.GetString("App.Host"), viper.GetString("App.Port"))); err != nil {
		log.Fatal("[Server start failed ] : ", err)
	}

}

func loadTranslation() {
	i18nBundle = i18n.NewBundle(language.Persian)
	i18nBundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	_, err := i18nBundle.LoadMessageFile("./internal/translation/active.fa.yaml")
	if err != nil {
		log.Fatal(err)
		return
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
