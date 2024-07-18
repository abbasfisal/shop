package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
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
	if err := db.Migrator().DropTable(
		&entities.User{},
		&entities.Category{},
		&entities.Product{},
		&entities.ProductImages{},
		&entities.Address{},
		&entities.Cart{},
		&entities.Order{},
	); err != nil {
		log.Fatal("[Drop tables] failed : ", err)
	}
	fmt.Println("[Drop] tables successfully")

	err := db.AutoMigrate(
		&entities.User{},
		&entities.Category{},
		&entities.Attribute{},
		&entities.AttributeValue{},
		&entities.Product{},
		&entities.ProductImages{},
		&entities.Address{},
		&entities.Cart{},
		&entities.Order{},
	)
	if err != nil {
		log.Fatal("[Migrate] table migration failed ", err)
		return
	}

	fmt.Println("[Create] tables successfully")

}
