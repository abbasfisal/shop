package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	App     App
	DB      DB
	Redis   Redis
	Upload  Upload
}
type App struct {
}
type DB struct {
}
type Redis struct {
}
type Upload struct {
}

func NewConfig() Config {
	err := godotenv.Load(viper.GetString("config.env_file"))
	if err != nil {
		log.Fatal("load .env file was failed : ", err)
	}

	return Config{}
}
