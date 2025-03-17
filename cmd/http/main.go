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
	"shop/cmd/job/worker"
	"shop/internal/database/mongodb"
	"shop/internal/database/mysql"
	"shop/internal/events"
	AdminRoutes "shop/internal/modules/admin/routes"
	PublicRoutes "shop/internal/modules/public/routes"
	"shop/internal/pkg/bootstrap"
	"shop/internal/pkg/cache"
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

	setupLog()

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	commands.Execute() // run cobra commands

	go worker.RunWorker(ctx, dependencies) // run Asynq worker
	go RunHttpServer(ctx, dependencies)    // run http server

	select {
	case <-ctx.Done():

		log.Println("[Shutting down] Closing database connections and external clients... after 5 second")

		<-time.After(time.Second * 5) // shut down after 5 second

		mysql.Close()
		mongodb.Disconnect()
		dependencies.AsynqClient.Close()

		log.Println("[Shutdown complete] All resources released. good by ;) ")

	}
}

func RunHttpServer(ctx context.Context, dependencies *bootstrap.Dependencies) {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	//logger middleware
	//r.Use(func(c *gin.Context) {
	//	//store in files only log
	//	start := time.Now()
	//
	//	c.Next()
	//
	//	log.Printf(
	//		"%s | %s | %s | %d | %s",
	//		time.Now().Format("2006-01-02 15:04:05"),
	//		c.Request.Method,
	//		c.Request.URL.Path,
	//		c.Writer.Status(),
	//		time.Since(start),
	//	)
	//})

	r.SetFuncMap(template.FuncMap{
		"stringToUint": util.StringToUint,
	})

	setupSessions(r)
	setupRoutes(ctx, r, dependencies)

	addr := fmt.Sprintf("%s:%s", viper.GetString("App.Host"), viper.GetString("App.Port"))
	log.Printf("[start server ]: http://%s\n", addr)
	if err := r.Run(addr); err != nil {
		logging.GlobalLog.FatalF("[Server start failed]: %v", err)
	}
}

// setupLog store logs in file and print in stdOut
func setupLog() {

	fileWriter := &lumberjack.Logger{
		Filename:   "./storage/logs/shop.log",
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

func setupRoutes(ctx context.Context, r *gin.Engine, dep *bootstrap.Dependencies) {

	r.LoadHTMLGlob("../../internal/**/**/**/*.html")
	r.Static("/uploads", "../../uploads")
	r.Static("/assets", "../../assets")
	r.StaticFile("/favicon.ico", "../../assets/shop/img/seller-logo.png")

	//eventManager dependencies : note: we have to prevent cycle import ,
	//									therefor we create another dependency struct
	eventManagerDep := events.EventManagerDep{
		AsynqClient: dep.AsynqClient,
		DB:          mysql.Get(),
		RedisClient: cache.NewRedisClient(),
		MongoClient: mongodb.Get(),
	}
	em := events.NewEventManager(&eventManagerDep)

	AdminRoutes.SetAdminRoutes(r, dep)
	PublicRoutes.SetPublic(r, dep, em)

	r.GET("/500", func(c *gin.Context) {
		c.HTML(http.StatusInternalServerError, "templates/html/errors/500", nil)
	})
}
