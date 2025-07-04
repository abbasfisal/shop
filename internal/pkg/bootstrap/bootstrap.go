package bootstrap

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/typesense/typesense-go/v3/typesense"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"shop/internal/database/mongodb"
	"shop/internal/database/mysql"
	"shop/internal/database/typesenceclient"
	"shop/internal/pkg/cache"
	"shop/internal/pkg/logging"
	"shop/internal/pkg/util"
	"sync"
)

var (
	once    sync.Once
	dep     *Dependencies
	initErr error
)

type Dependencies struct {
	I18nBundle      *i18n.Bundle
	AsynqClient     *asynq.Client
	DB              *gorm.DB
	RedisClient     *redis.Client
	MongoClient     *mongo.Client
	Storage         *util.Storage
	TypeSenceClient *typesense.Client
	Log             *logrus.Logger
}

func Initialize() (*Dependencies, error) {
	once.Do(func() {

		// load .env file
		if err := LoadConfig(); err != nil {
			initErr = err
			return
		}

		// config translation
		bundle, err := loadTranslation()
		if err != nil {
			initErr = err
			return
		}

		cache.InitRedisClient()   // redis connect
		mongodb.Connect()         // mongodb connect
		mysql.Connect()           // mysql connect
		typesenceclient.Connect() // initialize typesence

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
			Storage: util.NewStorage(os.Getenv("STORAGE_BUCKET_NAME"), os.Getenv("STORAGE_ENDPOINT_URL"),
				os.Getenv("STORAGE_ACEESS_KEY"), os.Getenv("STORAGE_SECRET_KEY")),
			//	EventManager: events.NewEventManager(&eventManagerDep),

			TypeSenceClient: typesenceclient.GetTClient(),
			Log:             logging.InitLogrus(),
		}

	})
	return dep, initErr
}

func LoadConfig() error {

	rootPath, _ := os.Getwd()
	envPath := filepath.Join(rootPath, ".env")

	if err := godotenv.Load(envPath); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	viper.AutomaticEnv() //read local environment automatically

	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	viper.AddConfigPath(filepath.Join(rootPath, "config"))
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("[x] viper reading config file was failed: %w", err)
	}

	return nil
}

func loadTranslation() (*i18n.Bundle, error) {
	i18nBundle := i18n.NewBundle(language.Persian)
	i18nBundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	rootPath, _ := os.Getwd()
	i18Path := filepath.Join(rootPath, "internal/translation/active.fa.yaml")

	if _, err := i18nBundle.LoadMessageFile(i18Path); err != nil {
		return nil, fmt.Errorf("[x] error loading translation file: %w", err)
	}

	return i18nBundle, nil
}

func initializeAsynqClient() (*asynq.Client, error) {
	opt := asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", viper.GetString("REDIS_DB"), viper.GetString("REDIS_PORT"))}
	return asynq.NewClient(opt), nil
}
