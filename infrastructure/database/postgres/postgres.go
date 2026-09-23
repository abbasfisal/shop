package postgres

import (
	"fmt"
	"log"
	"os"
	"strings"

	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Connect establishes the PostgreSQL connection (gorm).
func Connect() {
	// quote password so empty values don't break the keyword/value DSN parsing
	pw := strings.ReplaceAll(os.Getenv("POSTGRES_PASSWORD"), "'", "''")
	dsn := fmt.Sprintf(
		"host=%s user=%s password='%s' dbname=%s port=%s sslmode=%s TimeZone=UTC",
		env("POSTGRES_HOSTNAME", "127.0.0.1"),
		env("POSTGRES_USER", "postgres"),
		pw,
		env("POSTGRES_DB", "shop"),
		env("POSTGRES_PORT", "5432"),
		env("POSTGRES_SSLMODE", "disable"),
	)

	var err error
	db, err = gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("postgres error connection ", err)
	}

	if os.Getenv("APP_DEBUG") == "true" {
		db = db.Debug()
	}

	fmt.Println("\n[postgres] connected to postgresql database successfully ")
}

func Get() *gorm.DB {
	return db
}

func Close() error {
	connection, _ := db.DB()
	return connection.Close()
}
