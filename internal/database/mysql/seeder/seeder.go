package seeder

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

		Title:    "لباس و پوشاک",
		Slug:     "clothes",
		ParentID: nil,
		Image:    "",
		Status:   true,
		SubCategories: []entities.Category{
			{
				Model:    gorm.Model{},
				Title:    "لباس مردانه",
				Slug:     "men-clothes",
				ParentID: nil,
				Image:    "",
				Status:   true,
				SubCategories: []entities.Category{
					{
						Model:         gorm.Model{},
						Title:         "تی شرت",
						Slug:          "men-t-shirt",
						ParentID:      nil,
						Image:         "men-t-shirt.jpg",
						Status:        true,
						SubCategories: nil,
						Products: []entities.Product{
							{
								Model:         gorm.Model{},
								CategoryID:    0,
								Title:         "تیشرت مردانه مدل هایما",
								Slug:          "hima-men-t-shirt",
								Sku:           "sku-32932",
								Status:        true,
								OriginalPrice: 20000,
								SalePrice:     250000,
								Description:   "توضیحات تیشرت هایما",
								Category:      entities.Category{},
								ProductImages: nil,
							},
						},
					},
				},
				Products: nil,
			},
		},
		Products: nil,
	}

	var attributes = []entities.Attribute{
		{
			Model: gorm.Model{},
			Title: "سایز",

			AttributeValues: []entities.AttributeValue{
				{
					Model:          gorm.Model{},
					AttributeID:    0,
					AttributeTitle: "سایز",
					Value:          "s",
				},
				{
					Model:          gorm.Model{},
					AttributeID:    0,
					AttributeTitle: "سایز",
					Value:          "m",
				},
				{
					Model:          gorm.Model{},
					AttributeID:    0,
					AttributeTitle: "سایز",
					Value:          "l",
				},
				{
					Model:          gorm.Model{},
					AttributeID:    0,
					AttributeTitle: "سایز",
					Value:          "xl",
				},
				{
					Model:          gorm.Model{},
					AttributeID:    0,
					AttributeTitle: "سایز",
					Value:          "xxl",
				},
			},
		},
		{
			Model: gorm.Model{},
			Title: "رنگ",
			AttributeValues: []entities.AttributeValue{
				{
					Model:          gorm.Model{},
					AttributeID:    0,
					AttributeTitle: "رنگ",
					Value:          "آبی",
				},
				{
					Model:          gorm.Model{},
					AttributeID:    0,
					AttributeTitle: "رنگ",
					Value:          "قرمز",
				},
				{
					Model:          gorm.Model{},
					AttributeID:    0,
					AttributeTitle: "رنگ",
					Value:          "بنفش",
				},
			},
		},
	}

	db.Create(&user)
	db.Create(&category)
	db.Create(&attributes)
	fmt.Println("[Seed] tables successfully")
}
