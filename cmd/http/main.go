package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shop/cmd/commands"
	"shop/cmd/job/scheduler"
	"shop/cmd/job/worker"
	"shop/internal/database/mongodb"
	"shop/internal/database/mysql"
	"shop/internal/events"
	AdminRoutes "shop/internal/modules/admin/routes"
	PublicRoutes "shop/internal/modules/public/routes"
	"shop/internal/pkg/bootstrap"
	"shop/internal/pkg/logging"
	"shop/internal/pkg/util"
	"syscall"
	"time"
)

func main() {

	dependencies, err := bootstrap.Initialize()
	if err != nil {
		log.Fatal("[x] error initializing project :", err)
	}

	//eventManager dependencies : note: we have to prevent cycle import ,
	//									therefor we create another dependency struct
	eventManagerDep := events.EventManagerDep{
		AsynqClient: dependencies.AsynqClient,
		DB:          dependencies.DB,
		RedisClient: dependencies.RedisClient,
		MongoClient: dependencies.MongoClient,
	}
	em := events.NewEventManager(&eventManagerDep)

	setupLog()

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	commands.Execute() // run cobra commands

	go worker.RunWorker(ctx, dependencies, em)       // run Asynq worker
	go scheduler.RunScheduler(ctx, dependencies, em) // run Asynq Schedule
	go RunHttpServer(ctx, dependencies, em)          // run http server

	select {
	case <-ctx.Done():

		log.Println("[Shutting down] Closing database connections and external clients... after 5 second")

		<-time.After(time.Second) // shut down after 5 second

		mysql.Close()
		mongodb.Disconnect()
		dependencies.AsynqClient.Close()

		log.Println("[Shutdown complete] All resources released. good by ;) ")

	}
}

func RunHttpServer(ctx context.Context, dependencies *bootstrap.Dependencies, em *events.EventManager) {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.SetFuncMap(template.FuncMap{
		"stringToUint": util.StringToUint,
	})

	setupSessions(r)
	setupRoutes(ctx, r, dependencies, em)

	addr := fmt.Sprintf("%s:%s", viper.GetString("App.Host"), viper.GetString("App.Port"))
	log.Printf("[start server ]: http://%s\n", addr)
	if err := r.Run(addr); err != nil {
		logging.GlobalLog.FatalF("[Server start failed]: %v", err)
	}
}

// setupLog store logs in file and print in stdOut
func setupLog() {

	fileWriter := &lumberjack.Logger{
		Filename:   "../../storage/logs/shop.log",
		MaxSize:    10, //MB
		MaxAge:     10, //day
		MaxBackups: 5,
		LocalTime:  false,
		Compress:   true,
	}
	multi := io.MultiWriter(os.Stdout, fileWriter)
	log.SetOutput(multi)
}

func setupSessions(r *gin.Engine) {
	store := cookie.NewStore([]byte(viper.GetString("App.Key")))
	r.Use(sessions.Sessions("session", store))
}

func setupRoutes(ctx context.Context, r *gin.Engine, dep *bootstrap.Dependencies, em *events.EventManager) {

	r.LoadHTMLGlob("../../internal/**/**/**/*.html")
	r.Static("/uploads", "../../uploads")
	r.Static("/assets", "../../assets")
	r.StaticFile("/favicon.ico", "../../assets/shop/img/seller-logo.png")

	AdminRoutes.SetAdminRoutes(r, dep)
	PublicRoutes.SetPublic(r, dep, em)

	r.GET("/500", func(c *gin.Context) {
		c.HTML(http.StatusInternalServerError, "templates/html/errors/500", nil)
	})
}
