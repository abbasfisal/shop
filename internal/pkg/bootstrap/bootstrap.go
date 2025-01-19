package bootstrap

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"os"
	"shop/internal/database/mongodb"
	"shop/internal/database/mysql"
	"shop/internal/pkg/cache"
	"shop/internal/pkg/logging"
	"shop/internal/pkg/sms"
	"sync"
)

var (
	once    sync.Once
	dep     *Dependencies
	initErr error
)

type Dependencies struct {
	I18nBundle  *i18n.Bundle
	AsynqClient *asynq.Client
	DB          *gorm.DB
	RedisClient *redis.Client
	MongoClient *mongo.Client
}

func Initialize() (*Dependencies, error) {
	once.Do(func() {

		// load .env file
		if err := loadConfig(); err != nil {
			initErr = err
			return
		}

		// config translation
		bundle, err := loadTranslation()
		if err != nil {
			initErr = err
			return
		}

		cache.InitRedisClient() // redis connect
		mongodb.Connect()       // mongodb connect
		mysql.Connect()         // mysql connect

		// initialize logger
		initializeLogger()

		// initialize SMS
		initializeSmsService()

		// initialize Asynq
		asynqClient, err := initializeAsynqClient()
		if err != nil {
			initErr = err
			return
		}

		dep = &Dependencies{
			I18nBundle:  bundle,
			AsynqClient: asynqClient,
			DB:          mysql.Get(),
			RedisClient: cache.NewRedisClient(),
			MongoClient: mongodb.Get(),
			//	EventManager: events.NewEventManager(&eventManagerDep),
		}

	})
	return dep, initErr
}

func loadConfig() error {

	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("[x] failed to load .env file :%w", err)
	}

	viper.AutomaticEnv() //read local environment automatically

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config/")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("[x] viper reading config file was failed: %w", err)
	}

	return nil
}

func loadTranslation() (*i18n.Bundle, error) {
	i18nBundle := i18n.NewBundle(language.Persian)
	i18nBundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	if _, err := i18nBundle.LoadMessageFile("./internal/translation/active.fa.yaml"); err != nil {
		return nil, fmt.Errorf("[x] error loading translation file: %w", err)
	}

	return i18nBundle, nil
}

func initializeAsynqClient() (*asynq.Client, error) {
	opt := asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", viper.GetString("REDIS_DB"), viper.GetString("REDIS_PORT"))}
	return asynq.NewClient(opt), nil
}

func initializeSmsService() {
	kaveNegar := sms.NewKaveNegar(os.Getenv("KAVENEGAR_SECRETKEY"))
	sms.GetSMSManager().SetService(kaveNegar)
}

func initializeLogger() {
	logging.GlobalLog = logging.NewZapLogger()
}
