package commands

import (
	"database/sql"
	"fmt"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"shop/internal/database/mysql"
	"shop/internal/entities"
	"shop/internal/pkg/bootstrap"
)

var gorm bool
var up bool
var down bool
var steps int

// migrate command example:
//go run . migrate --up
//go run . migrate --down
//go run . migrate --down --steps=x

func init() {

	migrationCmd.Flags().BoolVar(&gorm, "gorm", false, "migrate with gorm using Models")

	migrationCmd.Flags().BoolVar(&up, "up", false, "Apply migration (upgrade)")
	migrationCmd.Flags().BoolVar(&down, "down", false, "Rollback migration (downgrade)")
	migrationCmd.Flags().IntVar(&steps, "steps", 1, "Number of migrations to rollback (default: 1 for --down)") // Added flag

	rootCmd.AddCommand(migrationCmd)
}

var migrationCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate Database Tables",
	Run: func(cmd *cobra.Command, args []string) {
		doMigrate()
	},
}

func doMigrate() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOSTNAME"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DB"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("sql open failed : ", err)
	}
	defer db.Close()

	var migrateType migrate.MigrationDirection
	maxSteps := 0

	//check flags
	if up {
		fmt.Println("Applying migration... (Upgrade)")
		migrateType = migrate.Up

	} else if down {
		fmt.Println("Rolling back migration... (Downgrade)")
		migrateType = migrate.Down
		maxSteps = steps

	} else if gorm {
		fmt.Println("migrate with gorm using Models")

		migrate1()

	} else {
		fmt.Println("Please specify a valid flag: -up or -down -gorm")
		os.Exit(1)

	}

	wd, _ := os.Getwd()

	migrations := &migrate.FileMigrationSource{
		Dir: filepath.Join(wd, "internal/database/mysql/migrations"),
	}

	n, err := migrate.ExecMax(db, "mysql", migrations, migrateType, maxSteps)
	if err != nil {
		log.Fatal("-- migration execute error:", err)
	}

	fmt.Printf("Migrations Done! , number of applies: %v\n", n)
}

func migratex() {
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
	//در مواقعی که اولین بار میخواهیم پروژه رو در محیط اجرا کنیم و تمام مایگریشن ها رو میخاهیم که اجرا بشه با این دستور میریم جلو
	//
	_, errs := bootstrap.Initialize()
	if errs != nil {
		log.Fatalln("err:", errs)
	}

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
