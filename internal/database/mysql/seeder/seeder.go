package seeder

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"shop/internal/database/mysql"
	"shop/internal/entities"
)

func Seed() {
	db := mysql.Get()

	hashPass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	var user = entities.User{
		FirstName:   "ali",
		LastName:    "mohammadi",
		PhoneNumber: "0935111111",
		Password:    string(hashPass),
		Type:        "admin",
	}

	var category = entities.Category{
		Title:  "Men Clothes",
		Slug:   "men-clothes",
		Image:  "categories/men.jpg",
		Status: true,
		Products: []entities.Product{
			{
				CategoryID: 0,
				Title:      "T shirt",
				Slug:       "t-shirt",
				Sku:        "sku1000",
				Status:     true,
				//Quantity:      100,
				OriginalPrice: 600000,
				SalePrice:     500000,
				Description:   "description goes here",
			},
			{
				CategoryID: 0,
				Title:      "Belt",
				Slug:       "belt",
				Sku:        "sku2000",
				Status:     true,
				//Quantity:      50,
				OriginalPrice: 200000,
				SalePrice:     100000,
				Description:   "description goes here",
			},
		},
	}

	//address
	var address = entities.Address{
		UserID:      1,
		Floor:       1,
		Number:      5,
		Phase:       "A",
		Block:       "B",
		Description: "address description goes here",
	}

	var cart1 = entities.Cart{
		UserID:    1,
		ProductID: 1,
		Count:     5,
		Status:    0, //not paid
	}
	var cart2 = entities.Cart{
		UserID:    1,
		ProductID: 2,
		Count:     4,
		Status:    0, //not paid
	}

	var order1 = entities.Order{
		UserID:         1,
		TrackingNumber: "879821",
		ProductID:      1,
		Count:          5,
		Status:         1, //successfully paid
	}

	var order2 = entities.Order{
		UserID:         1,
		TrackingNumber: "879821",
		ProductID:      2,
		Count:          4,
		Status:         1, //successfully paid
	}

	db.Create(&user)
	db.Create(&category)
	db.Create(&address)

	db.Create(&cart1)
	db.Create(&cart2)

	db.Create(&order1)
	db.Create(&order2)

	fmt.Println("[Seed] tables successfully")
}
