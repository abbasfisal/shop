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

	//user
	var user = entities.User{
		FirstName:   "ali",
		LastName:    "mohammadi",
		PhoneNumber: "0935111111",
		Password:    string(hashPass),
		Type:        "admin",
	}

	//brand
	var brand = entities.Brand{
		Model:   gorm.Model{},
		Title:   "zara",
		Slug:    "zara",
		Image:   "",
		Product: nil,
	}

	//category
	var category = entities.Category{
		Model:    gorm.Model{},
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
					},
				},
				Products: nil,
			},
		},
		Products: nil,
	}

	//attribute -- attribute-values
	var attributesAndValues = []entities.Attribute{
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

	//attribute-value

	//product
	var product = entities.Product{
		Model:         gorm.Model{},
		BrandID:       1,
		CategoryID:    3,
		Title:         "تیشرت مردانه مدل هایما",
		Slug:          "hima-men-t-shirt",
		Sku:           "sku-32932",
		Status:        true,
		OriginalPrice: 20000,
		SalePrice:     250000,
		Description:   "توضیحات تیشرت هایما",
	}

	//product-attribute
	var productAttribute = []entities.ProductAttribute{
		{
			Model:               gorm.Model{},
			ProductID:           1,
			AttributeID:         1, //size
			AttributeTitle:      "سایز",
			AttributeValueID:    1, //small
			AttributeValueTitle: "s",
		},
		{
			Model:               gorm.Model{},
			ProductID:           1,
			AttributeID:         1, //size
			AttributeTitle:      "سایز",
			AttributeValueID:    2, //medium
			AttributeValueTitle: "m",
		},
		{
			Model:               gorm.Model{},
			ProductID:           1,
			AttributeID:         2, //color
			AttributeTitle:      "رنگ",
			AttributeValueID:    6, //ابی
			AttributeValueTitle: "ابی",
		},
		{
			Model:               gorm.Model{},
			ProductID:           1,
			AttributeID:         2, //color
			AttributeTitle:      "رنگ",
			AttributeValueID:    7, //قرمز
			AttributeValueTitle: "قرمز",
		},
	}

	//product-inventory
	var productInventory = []entities.ProductInventory{
		{
			Model:     gorm.Model{},
			ProductID: 1,
			Quantity:  20,
		},
		{
			Model:     gorm.Model{},
			ProductID: 1,
			Quantity:  10,
		},
	}

	//product-inventory-attribute
	var productInventoryAttribute = []entities.ProductInventoryAttribute{
		//1 : size : small
		//2 : size : medium
		//3 : color: blue
		//4 : color: red

		//small-blue = qty =20
		{
			Model:              gorm.Model{},
			ProductID:          1,
			ProductInventoryID: 1, //qty=20
			ProductAttributeID: 1, //small
		},
		{
			Model:              gorm.Model{},
			ProductID:          1,
			ProductInventoryID: 1, //qty=20
			ProductAttributeID: 3, //blue
		},

		//medium-red = qty =10
		{
			Model:              gorm.Model{},
			ProductID:          1,
			ProductInventoryID: 1, //qty=10
			ProductAttributeID: 2, //medium
		},
		{
			Model:              gorm.Model{},
			ProductID:          1,
			ProductInventoryID: 1, //qty=10
			ProductAttributeID: 4, //red
		},
	}
	db.Create(&user)
	db.Create(&category)
	db.Create(&brand)
	db.Create(&attributesAndValues)
	db.Create(&product)
	db.Create(&productAttribute)
	db.Create(&productInventory)
	db.Create(&productInventoryAttribute)

	fmt.Println("[Seed] tables successfully")
}
