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
	var brand = []entities.Brand{
		{
			Title: "zara",
			Slug:  "zara",
		},
		//perfume brands
		{
			Title: "Bailando",
			Slug:  "bailando",
		},
		{
			Title: "Ballerina",
			Slug:  "ballerina",
			Image: "",
		},
		{
			Title: "woody sence",
			Slug:  "woody-sence",
		},
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
	var products = []entities.Product{
		//women perfume
		{
			CategoryID:    49, //عطر و ادکلن زنانه
			BrandID:       3,  //Ballerina
			Title:         "ادو پرفیوم زنانه بالرینا مدل گود گرل Good Girl حجم 90 میلی لیتر",
			Slug:          "ادو-پرفیوم-زنانه-بالرینا-مدل-گود-گرل-good-girl",
			Sku:           "sku1000",
			Status:        true,
			OriginalPrice: 822_000,
			SalePrice:     349_000,
			Description:   "ادو پرفیوم زنانه بالرینا مدل Good Girl عطری است که با رایحه ی منحصر به فرد خود به یکی از محبوب ترین عطرهای زنانه در دنیای عطر و ادکلن تبدیل شده است. این عطر مناسب خانم هایی است که به دنبال رایحه ای جذاب، ماندگار و خاص هستند.",
			ProductImages: []entities.ProductImages{
				{
					Path: "2024/09/27/1.webp",
				},
				{
					Path: "2024/09/27/2.webp",
				},
				{
					Path: "2024/09/27/3.webp",
				},
			},
			ProductInventories: []entities.ProductInventory{
				{
					Quantity: 150,
				},
			},
			Features: []entities.Feature{
				{
					Title: "نوع عطر",
					Value: "ادو پرفیوم",
				},
				{
					Title: "حجم",
					Value: "90 میلی لیتر",
				},
				{
					Title: "مناسب برای",
					Value: "خانم ها",
				},
				{
					Title: "نت های آغازین",
					Value: "بادام، قهوه، ترنج، ليمو",
				},
				{
					Title: "نت های میانی",
					Value: "گل سرخ، ياس سامباك، شكوفه پرتقال، زنبق و رز بلغاري",
				},
				{
					Title: "نت های پایانی",
					Value: "لوبياي تونكا\n غلاف كاكائو\nوانيل\nپرالين\nچوب صندل\nمشك\nعنبر\nنعناع هندي و سدر\nچوب كشمير\nدارچين\n\n",
				},
				{
					Title: "نوع رایحه",
					Value: "گرم و تلخ",
				},
				{
					Title: "فصل",
					Value: "پاییز و زمستان",
				},
			},
		},
		{
			CategoryID:    49, //عطر و ادکلن زنانه
			BrandID:       3,  //Ballerina
			Title:         "ادو پرفیوم زنانه بالرینا مدل پویزن Poisson حجم 100 میلی لیتر",
			Slug:          "ادو-پرفیوم-زنانه-بالرینا-مدل-پویزن-poisson",
			Sku:           "sku1001",
			Status:        true,
			OriginalPrice: 780_000,
			SalePrice:     349_000,
			Description:   "ادو پرفیوم زنانه بالرینا مدل پویزن Poisson عطری است زنانه با رایحه ای شیرین و گرم که مکمل شخصیت زنانه است و به شما احساس منحصر به فرد و جذاب می دهد. با بسته‌بندی و طراحی لوکس شیشه، این عطر بهترین کیفیت را در اختیار شما قرار می‌دهد.",
			ProductImages: []entities.ProductImages{
				{
					Path: "2024/09/27/11.webp",
				},
				{
					Path: "2024/09/27/22.webp",
				},
				{
					Path: "2024/09/27/33.webp",
				},
			},
			ProductInventories: []entities.ProductInventory{
				{
					Quantity: 200,
				},
			},
			Features: []entities.Feature{
				{
					Title: "نوع محصول (غلظت)",
					Value: "ادو پرفیوم",
				},
				{
					Title: "حجم",
					Value: "100 میلی لیتر",
				},
				{
					Title: "مناسب برای",
					Value: "خانم ها",
				},
				{
					Title: "نت های آغازین",
					Value: "ترنج و ماندارین",
				},
				{
					Title: "نت های میانی",
					Value: "گل رز",
				},
				{
					Title: "نت های پایانی",
					Value: "عنبر\n\n, \nنعنا هندی\n\n, \nوانيل",
				},
				{
					Title: "نوع رایحه",
					Value: "گرم و شیرین",
				},
				{
					Title: "فصل",
					Value: "پاییز و زمستان",
				},
				{
					Title: "گروه بویایی",
					Value: "چوبی چایپر",
				},
			},
		},
		{
			CategoryID:    49, //عطر و ادکلن زنانه
			BrandID:       3,  //Ballerina
			Title:         "ادو پرفیوم زنانه بایلندو مدل اکلت Eclatto حجم 100 میلی لیتر",
			Slug:          "ادو-پرفیوم-زنانه-بایلندو-مدل-اکلت-eclatto",
			Sku:           "sku1002",
			Status:        true,
			OriginalPrice: 815_000,
			SalePrice:     477_600,
			Description:   "ادو پرفیوم زنانه بایلندو مدل d’ Eclatto قصیده ای فریبنده برای ظرافت زنانگی است،‌ جاییکه ترکیب مست کننده میوه ها، لمس مخملی گل پائونیا، و با حضور باشکوه سرو گرد هم می آیند.تا نقش و نگار طلسم کننده ای از جذابیت و اعتماد به نفس را بیافریند.",
			ProductImages: []entities.ProductImages{
				{
					Path: "2024/09/27/111.webp",
				},
				{
					Path: "2024/09/27/222.webp",
				},
				{
					Path: "2024/09/27/333.webp",
				},
			},
			ProductInventories: []entities.ProductInventory{
				{
					Quantity: 200,
				},
			},
			Features: []entities.Feature{
				{
					Title: "نوع محصول (غلظت)",
					Value: "ادو پرفیوم",
				},
				{
					Title: "حجم",
					Value: "100 میلی لیتر",
				},
				{
					Title: "مناسب برای",
					Value: "خانم ها",
				},
				{
					Title: "نت های آغازین",
					Value: "درخت سرو\n\n, \nگل پائونیا",
				},
				{
					Title: "نت های میانی",
					Value: "پونه\n\n, \nمشک\n\n, \nويستريا",
				},
				{
					Title: "نت های پایانی",
					Value: "اسمانتوس\n\n, \nسدر\n\n, \nعنبر",
				},
				{
					Title: "نوع رایحه",
					Value: "ملایم و شیرین",
				},
				{
					Title: "فصل",
					Value: "بهار",
				},
			},
		},
	}

	db.Create(&user)
	db.Create(&category)
	db.Create(&brand)
	db.Create(&attributesAndValues)
	db.Create(&products)
	//db.Create(&productAttribute)
	//db.Create(&productInventory)
	//db.Create(&productInventoryAttribute)

	fmt.Println("\\\\\\\\\\\\\\  ~~~~[Seed] tables successfully~~~~ \\\\\\\\\\\\\\")
}
