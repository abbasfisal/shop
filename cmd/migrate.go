package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"shop/infrastructure/database/postgres"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(migrateStatusCmd)
	rootCmd.AddCommand(migrateRollbackCmd)
	rootCmd.AddCommand(migrateResetCmd)
	rootCmd.AddCommand(makeMigrationCmd)
}

func openMigrateDB() *sql.DB {
	postgres.Connect()
	sqlDB, err := postgres.Get().DB()
	if err != nil {
		log.Fatal("postgres sql handle failed: ", err)
	}
	return sqlDB
}

func migrationsDir() string {
	if dir := os.Getenv("MIGRATIONS_DIR"); dir != "" {
		return dir
	}
	wd, _ := os.Getwd()
	return filepath.Join(wd, "migrations")
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply all pending migrations (goose up)",
	Run: func(cmd *cobra.Command, args []string) {
		db := openMigrateDB()
		defer db.Close()

		if err := goose.SetDialect("postgres"); err != nil {
			log.Fatal(err)
		}
		if err := goose.Up(db, migrationsDir()); err != nil {
			log.Fatal("-- migration up error: ", err)
		}
		fmt.Println("Migrations applied successfully")
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "migrate:status",
	Short: "Show migration status",
	Run: func(cmd *cobra.Command, args []string) {
		db := openMigrateDB()
		defer db.Close()

		if err := goose.SetDialect("postgres"); err != nil {
			log.Fatal(err)
		}
		if err := goose.Status(db, migrationsDir()); err != nil {
			log.Fatal("-- migration status error: ", err)
		}
	},
}

var migrateRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the last migration (goose down)",
	Run: func(cmd *cobra.Command, args []string) {
		db := openMigrateDB()
		defer db.Close()

		if err := goose.SetDialect("postgres"); err != nil {
			log.Fatal(err)
		}
		if err := goose.Down(db, migrationsDir()); err != nil {
			log.Fatal("-- migration rollback error: ", err)
		}
		fmt.Println("Migration rolled back successfully")
	},
}

var migrateResetCmd = &cobra.Command{
	Use:   "migrate:reset",
	Short: "Rollback all migrations then apply them again",
	Run: func(cmd *cobra.Command, args []string) {
		db := openMigrateDB()
		defer db.Close()

		if err := goose.SetDialect("postgres"); err != nil {
			log.Fatal(err)
		}
		if err := goose.Reset(db, migrationsDir()); err != nil {
			log.Fatal("-- migration reset error: ", err)
		}
		fmt.Println("Migrations reset successfully")
	},
}

var migrationName string

var makeMigrationCmd = &cobra.Command{
	Use:   "make:migration",
	Short: "Create a new migration file (e.g. make:migration NAME=add_users_table)",
	Run: func(cmd *cobra.Command, args []string) {
		if migrationName == "" && len(args) > 0 {
			migrationName = args[0]
		}
		if migrationName == "" {
			log.Fatal("please provide a migration name: make:migration NAME=add_xxx")
		}

		stamp := time.Now().Format("20060102150405")
		fileName := fmt.Sprintf("%s_%s.sql", stamp, strings.ToLower(migrationName))
		path := filepath.Join(migrationsDir(), fileName)

		content := "-- +goose Up\n-- +goose Down\n"
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			log.Fatal("create migration failed: ", err)
		}
		fmt.Println("created:", path)
	},
}

func init() {
	makeMigrationCmd.Flags().StringVar(&migrationName, "name", "", "migration name")
}
