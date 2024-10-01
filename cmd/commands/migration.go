package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"shop/internal/database/mysql"
	"shop/internal/entities"
)

func init() {
	rootCmd.AddCommand(migrationCmd)
}

var migrationCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate Database Tables",
	Run: func(cmd *cobra.Command, args []string) {
		migrate()
	},
}

func migrate() {

	db := mysql.Get()
	db.Exec("SET foreign_key_checks = 0")

	if err := db.Migrator().DropTable(
		&entities.Cart{},
		&entities.User{},
		&entities.Customer{},
		&entities.OTP{},
		&entities.Session{},
		&entities.Brand{},
		&entities.Category{},
		&entities.ProductInventoryAttribute{},
		&entities.ProductInventory{},
		&entities.Product{},
		&entities.Feature{},
		&entities.ProductAttribute{},
		&entities.ProductImages{},
		&entities.AttributeValue{},
		&entities.Attribute{},
		&entities.Address{},
		&entities.Order{},
	); err != nil {
		log.Fatal("[Drop tables] failed : ", err)
	}
	fmt.Println("[Drop] tables successfully")

	err := db.AutoMigrate(
		&entities.User{},
		&entities.Customer{},
		&entities.OTP{},
		&entities.Session{},
		&entities.Brand{},
		&entities.Category{},
		&entities.Feature{},
		&entities.ProductInventoryAttribute{},
		&entities.Attribute{},
		&entities.AttributeValue{},
		&entities.Product{},
		&entities.ProductAttribute{},
		&entities.ProductImages{},
		&entities.ProductInventory{},
		&entities.Address{},
		&entities.Cart{},
		&entities.Order{},
	)
	db.Exec("SET foreign_key_checks = 1")

	if err != nil {
		log.Fatal("~~~~ [Migrate] Tables Migration Failed : ", err)
		return
	}

	fmt.Println("~~~~ [Migrate] Tables Migration Success")

	os.Exit(1)
}
