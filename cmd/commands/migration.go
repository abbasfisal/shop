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
		&entities.Payment{},
		&entities.OrderItem{},
		&entities.Order{},
		&entities.CartItem{},
		&entities.Cart{},
	); err != nil {
		log.Fatal("[Drop tables] failed : ", err)
	}
	fmt.Println("[Drop] tables successfully")

	err := db.AutoMigrate(
		&entities.Order{},
		&entities.Payment{},
		&entities.OrderItem{},
		&entities.Cart{},
		&entities.CartItem{},
	)
	db.Exec("SET foreign_key_checks = 1")

	if err != nil {
		log.Fatal("~~~~ [Migrate] Tables Migration Failed : ", err)
		return
	}

	fmt.Println("~~~~ [Migrate] Tables Migration Success")

	os.Exit(1)
}
func migrate1() {

	db := mysql.Get()
	db.Exec("SET foreign_key_checks = 0")

	if err := db.Migrator().DropTable(
		&entities.OrderItem{},
		&entities.Payment{},
		&entities.Order{},
		&entities.CartItem{},
		&entities.Cart{},
		&entities.Address{},
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
	); err != nil {
		log.Fatal("[Drop tables] failed : ", err)
	}
	fmt.Println("[Drop] tables successfully")

	err := db.AutoMigrate(
		&entities.User{},
		&entities.Customer{},
		&entities.Address{},
		&entities.Brand{},
		&entities.Category{},
		&entities.Attribute{},
		&entities.AttributeValue{},
		&entities.Product{},
		&entities.ProductImages{},
		&entities.Feature{},
		&entities.ProductAttribute{},
		&entities.ProductInventory{},
		&entities.ProductInventoryAttribute{},
		&entities.Cart{},
		&entities.CartItem{},
		&entities.Order{},
		&entities.OrderItem{},
		&entities.Payment{},
		&entities.OTP{},
		&entities.Session{},
	)
	db.Exec("SET foreign_key_checks = 1")

	if err != nil {
		log.Fatal("~~~~ [Migrate] Tables Migration Failed : ", err)
		return
	}

	fmt.Println("~~~~ [Migrate] Tables Migration Success")

	os.Exit(1)
}
