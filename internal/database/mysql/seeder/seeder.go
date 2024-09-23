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
	id1 := uint(1)
	id2 := uint(2)

	category := []entities.Category{
		//apparel
		{
			Priority: &id1,
			Title:    "مد و پوشاک",
			Slug:     "apparel",
			ParentID: nil,
			Status:   true,
			SubCategories: []entities.Category{
				//men clothing
				{
					Priority: nil,
					Title:    "لباس مردانه",
					Slug:     "category-men-clothing",
					ParentID: nil,
					Image:    "",
					Status:   true,
					SubCategories: []entities.Category{
						{
							Priority: nil,
							Title:    "تی شرت مردانه",
							Slug:     "category-men-tee-shirts",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "پیراهن مردانه",
							Slug:     "category-men-shirts",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "لباس زیر مردانه",
							Slug:     "category-men-underwear",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "شلوار مردانه",
							Slug:     "category-men-trousers-jumpsuits",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "سویشرت مردانه",
							Slug:     "category-men-sweatshirts",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "ژاکت و پلیور مردانه",
							Slug:     "category-men-knitwear",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "هودی مردانه",
							Slug:     "category-men-hoodies",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
					},
				},

				//men accessories
				{
					Priority: nil,
					Title:    "اکسسوری مردانه",
					Slug:     "category-men-accessories",
					ParentID: nil,
					Image:    "",
					Status:   true,
					SubCategories: []entities.Category{
						{
							Priority:      nil,
							Title:         "عینک مردانه",
							Slug:          "category-men-eyewear",
							ParentID:      nil,
							Image:         "",
							Status:        true,
							SubCategories: nil,
							Products:      nil,
						},
						{
							Priority:      nil,
							Title:         "کمربند و ساسبند مردانه",
							Slug:          "category-men-belt",
							ParentID:      nil,
							Image:         "",
							Status:        true,
							SubCategories: nil,
							Products:      nil,
						},
						{
							Priority:      nil,
							Title:         "کلاه مردانه",
							Slug:          "category-men-headwear",
							ParentID:      nil,
							Image:         "",
							Status:        true,
							SubCategories: nil,
							Products:      nil,
						},
						{
							Priority:      nil,
							Title:         "کراوات و پاپیون مردانه",
							Slug:          "category-men-ties",
							ParentID:      nil,
							Image:         "",
							Status:        true,
							SubCategories: nil,
							Products:      nil,
						},
					},
					Products: nil,
				},

				//women home wear
				{
					Priority: nil,
					Title:    "لباس خواب و راحتی زنانه",
					Slug:     "category-women-home-wear",
					ParentID: nil,
					Image:    "",
					Status:   true,
					SubCategories: []entities.Category{
						{
							Priority: nil,
							Title:    "لباس خواب و راحتی زنانه",
							Slug:     "category-women-home-wear-set",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
					},
					Products: nil,
				},

				//women shirts
				{
					Priority: nil,
					Title:    "بلوز و شومیز زنانه",
					Slug:     "category-women-shirts",
					ParentID: nil,
					Image:    "",
					Status:   true,
					SubCategories: []entities.Category{
						{
							Priority: nil,
							Title:    "بادی زنانه",
							Slug:     "category-women-body-suits",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "ژاکت و پلیور زنانه",
							Slug:     "category-women-knitwear",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "تاپ و نیم تنه زنانه",
							Slug:     "category-women-tops-and-croptops",
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "تونیک زنانه",
							Slug:     "category-women-tunic",
							ParentID: nil,
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "تیشرت زنانه",
							Slug:     "category-women-tee-shirts",
							ParentID: nil,
							Status:   true,
						},
					},
				},

				//woman underwear
				{
					Priority: nil,
					Title:    "لباس زیر زنانه",
					Slug:     "category-women-underwear",
					ParentID: nil,
					Image:    "",
					Status:   true,
					SubCategories: []entities.Category{
						{
							Priority: nil,
							Title:    "شورت زنانه",
							Slug:     "category-women-underwear-bottom",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "سوتین زنانه",
							Slug:     "category-women-bra",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
					},
				},

				//woman accessories
				{
					Priority: nil,
					Title:    "اکسسوری زنانه",
					Slug:     "category-women-accessories",
					ParentID: nil,
					Image:    "",
					Status:   true,
					SubCategories: []entities.Category{
						{
							Priority: nil,
							Title:    "ساعت زنانه",
							Slug:     "category-women-watches",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "عینک زنانه",
							Slug:     "category-women-eyewear",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "کمربند زنانه",
							Slug:     "category-women-belts",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
					},
				},

				//girl clothes
				{
					Priority: nil,
					Title:    "لباس دخترانه",
					Slug:     "category-girls-clothes",
					ParentID: nil,
					Image:    "",
					Status:   true,
					SubCategories: []entities.Category{
						{
							Priority: nil,
							Title:    "بلوز و شومیز",
							Slug:     "category-girls-shirts",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
					},
				},

				//girl underwear
				{
					Priority: nil,
					Title:    "لباس زیر دخترانه",
					Slug:     "category-girls-underwear",
					ParentID: nil,
					Image:    "",
					Status:   true,
					SubCategories: []entities.Category{
						{
							Priority: nil,
							Title:    "شورت دخترانه",
							Slug:     "category-girls-underwear-bottom",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "سوتین دخترانه",
							Slug:     "category-girls-bra",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
					},
				},

				//girl accessories
				{
					Priority: nil,
					Title:    "اکسسوری دخترانه",
					Slug:     "category-girls-accessories",
					ParentID: nil,
					Image:    "",
					Status:   true,
					SubCategories: []entities.Category{
						{
							Priority: nil,
							Title:    "شال و روسری دخترانه",
							Slug:     "category-girls-scarves",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
					},
				},

				//boys clothes
				{
					Priority: nil,
					Title:    "لباس پسرانه",
					Slug:     "category-boys-clothes",
					ParentID: nil,
					Image:    "",
					Status:   true,
					SubCategories: []entities.Category{
						{
							Priority: nil,
							Title:    "تیشرت پسرانه",
							Slug:     "category-boys-tee-shirts",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "شلوار پسرانه",
							Slug:     "category-boys-trousers",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "سویشرت و هودی پسرانه",
							Slug:     "category-boys-hoodies",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
					},
				},

				//boys underwear
				{
					Priority: nil,
					Title:    "لباس زیر پسرانه",
					Slug:     "category-boys-underwear",
					ParentID: nil,
					Image:    "",
					Status:   true,
					SubCategories: []entities.Category{
						{
							Priority: nil,
							Title:    "شورت پسرانه",
							Slug:     "category-boys-underwear-bottom",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
					},
				},

				//mom fit
				{
					Priority: nil,
					Title:    "شلوار مام فیت",
					Slug:     "mom-fit-pants",
					ParentID: nil,
					Image:    "",
					Status:   true,
				},

				//flip-flop
				{
					Priority: nil,
					Title:    "شلوار دمپا",
					Slug:     "flip-flop-pants",
					ParentID: nil,
					Image:    "",
					Status:   true,
				},

				//baggy pant
				{
					Priority: nil,
					Title:    "شلوار بگ",
					Slug:     "baggy-pants",
					ParentID: nil,
					Image:    "",
					Status:   true,
				},

				//mom pant
				{
					Priority: nil,
					Title:    "شلوار مام",
					Slug:     "mom-pants",
					ParentID: nil,
					Image:    "",
					Status:   true,
				},
			},
		},

		//personal appliance
		{
			Priority: &id2,
			Title:    "آرایشی و بهداشتی",
			Slug:     "category-personal-appliance",
			ParentID: nil,
			Image:    "",
			Status:   true,
			SubCategories: []entities.Category{
				{
					Priority: nil,
					Title:    "عطر و ادکلن",
					Slug:     "category-perfume-all",
					ParentID: nil,
					Image:    "",
					Status:   true,
					SubCategories: []entities.Category{
						{
							Priority: nil,
							Title:    "عطر و ادکلن زنانه",
							Slug:     "category-women-perfume",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "عطر و ادکلن مردانه",
							Slug:     "category-men-perfume",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "بادی اسپلش",
							Slug:     "category-body-splash",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
						{
							Priority: nil,
							Title:    "عطر جیبی",
							Slug:     "category-pocket-perfume",
							ParentID: nil,
							Image:    "",
							Status:   true,
						},
					},
					Products: nil,
				},
			},
			Products: nil,
		},
	}
	//old category
	//var category = entities.Category{
	//	Model:    gorm.Model{},
	//	Title:    "لباس و پوشاک",
	//	Slug:     "clothes",
	//	ParentID: nil,
	//	Image:    "",
	//	Status:   true,
	//	SubCategories: []entities.Category{
	//		{
	//			Model:    gorm.Model{},
	//			Title:    "لباس مردانه",
	//			Slug:     "men-clothes",
	//			ParentID: nil,
	//			Image:    "",
	//			Status:   true,
	//			SubCategories: []entities.Category{
	//				{
	//					Model:         gorm.Model{},
	//					Title:         "تی شرت",
	//					Slug:          "men-t-shirt",
	//					ParentID:      nil,
	//					Image:         "men-t-shirt.jpg",
	//					Status:        true,
	//					SubCategories: nil,
	//				},
	//			},
	//			Products: nil,
	//		},
	//	},
	//	Products: nil,
	//}

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

	//my product
	//attribute-value
	//
	////product
	//var product = entities.Product{
	//	Model:         gorm.Model{},
	//	BrandID:       1,
	//	CategoryID:    3,
	//	Title:         "تیشرت مردانه مدل هایما",
	//	Slug:          "hima-men-t-shirt",
	//	Sku:           "sku-32932",
	//	Status:        true,
	//	OriginalPrice: 20000,
	//	SalePrice:     250000,
	//	Description:   "توضیحات تیشرت هایما",
	//}
	//
	////product-attribute
	//var productAttribute = []entities.ProductAttribute{
	//	{
	//		Model:               gorm.Model{},
	//		ProductID:           1,
	//		AttributeID:         1, //size
	//		AttributeTitle:      "سایز",
	//		AttributeValueID:    1, //small
	//		AttributeValueTitle: "s",
	//	},
	//	{
	//		Model:               gorm.Model{},
	//		ProductID:           1,
	//		AttributeID:         1, //size
	//		AttributeTitle:      "سایز",
	//		AttributeValueID:    2, //medium
	//		AttributeValueTitle: "m",
	//	},
	//	{
	//		Model:               gorm.Model{},
	//		ProductID:           1,
	//		AttributeID:         2, //color
	//		AttributeTitle:      "رنگ",
	//		AttributeValueID:    6, //ابی
	//		AttributeValueTitle: "ابی",
	//	},
	//	{
	//		Model:               gorm.Model{},
	//		ProductID:           1,
	//		AttributeID:         2, //color
	//		AttributeTitle:      "رنگ",
	//		AttributeValueID:    7, //قرمز
	//		AttributeValueTitle: "قرمز",
	//	},
	//}
	//
	////product-inventory
	//var productInventory = []entities.ProductInventory{
	//	{
	//		Model:     gorm.Model{},
	//		ProductID: 1,
	//		Quantity:  20,
	//	},
	//	{
	//		Model:     gorm.Model{},
	//		ProductID: 1,
	//		Quantity:  10,
	//	},
	//}
	//
	////product-inventory-attribute
	//var productInventoryAttribute = []entities.ProductInventoryAttribute{
	//	//1 : size : small
	//	//2 : size : medium
	//	//3 : color: blue
	//	//4 : color: red
	//
	//	//small-blue = qty =20
	//	{
	//		Model:              gorm.Model{},
	//		ProductID:          1,
	//		ProductInventoryID: 1, //qty=20
	//		ProductAttributeID: 1, //small
	//	},
	//	{
	//		Model:              gorm.Model{},
	//		ProductID:          1,
	//		ProductInventoryID: 1, //qty=20
	//		ProductAttributeID: 3, //blue
	//	},
	//
	//	//medium-red = qty =10
	//	{
	//		Model:              gorm.Model{},
	//		ProductID:          1,
	//		ProductInventoryID: 1, //qty=10
	//		ProductAttributeID: 2, //medium
	//	},
	//	{
	//		Model:              gorm.Model{},
	//		ProductID:          1,
	//		ProductInventoryID: 1, //qty=10
	//		ProductAttributeID: 4, //red
	//	},
	//}
	//

	//end my product

	db.Create(&user)
	db.Create(&category)
	db.Create(&brand)
	db.Create(&attributesAndValues)
	//db.Create(&product)
	//db.Create(&productAttribute)
	//db.Create(&productInventory)
	//db.Create(&productInventoryAttribute)

	fmt.Println("[Seed] tables successfully")
}
