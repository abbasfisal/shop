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
								ProductImage:  nil,
								//ProductAttribute: []entities.ProductAttribute{
								//	{
								//		Model:               gorm.Model{},
								//		ProductID:           0,
								//		AttributeID:         1,
								//		AttributeTitle:      "سایز",
								//		AttributeValueID:    1,
								//		AttributeValueTitle: "s",
								//		Attribute:           entities.Attribute{},
								//		AttributeValue:      entities.AttributeValue{},
								//	},
								//	{
								//		Model:               gorm.Model{},
								//		ProductID:           0,
								//		AttributeID:         1,
								//		AttributeTitle:      "سایز",
								//		AttributeValueID:    2,
								//		AttributeValueTitle: "m",
								//		Attribute:           entities.Attribute{},
								//		AttributeValue:      entities.AttributeValue{},
								//	},
								//	{
								//		Model:               gorm.Model{},
								//		ProductID:           0,
								//		AttributeID:         1,
								//		AttributeTitle:      "سایز",
								//		AttributeValueID:    3,
								//		AttributeValueTitle: "l",
								//		Attribute:           entities.Attribute{},
								//		AttributeValue:      entities.AttributeValue{},
								//	},
								//},
							},
						},
						Attribute: nil,
					},
				},
				Products:  nil,
				Attribute: nil,
			},
		},
		Products: nil,
		Attribute: []entities.Attribute{
			{
				Model:      gorm.Model{},
				CategoryID: 0,
				Title:      "سایز",
				Category:   entities.Category{},
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
				Model:      gorm.Model{},
				CategoryID: 0,
				Title:      "رنگ",
				Category:   entities.Category{},
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
		},
	}

	db.Create(&user)
	db.Create(&category)

	fmt.Println("[Seed] tables successfully")
}
