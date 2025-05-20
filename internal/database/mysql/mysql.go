package mysql

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB

func Connect() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOSTNAME"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DB"),
	)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("mysql error connection ", err)
	}

	if os.Getenv("APP_DEBUG") == "true" {
		db = db.Debug()
	}

	fmt.Println("\n[mysql] connected to mysql database successfully ")
}

func Get() *gorm.DB {
	return db
}

func Close() error {
	connection, _ := db.DB()
	err := connection.Close()
	if err != nil {
		return err
	}
	return nil
}
