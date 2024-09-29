package seeder

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"shop/internal/database/mysql"
	"shop/internal/entities"
	"strconv"
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
	var brand = fakeBrands()

	//category
	category := fakeCategories()

	//attribute -- attribute-values
	var attributesAndValues = fakeAttributeAndValues()

	//products include :
	//3 perfume
	//1 men trouser
	var products = fakeProducts()

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

func fakeProducts() []entities.Product {
	var p []entities.Product
	for i := 1; i <= 200; i++ {
		p = append(p, entities.Product{
			CategoryID:    49, //عطر و ادکلن زنانه
			BrandID:       3,  //Ballerina
			Title:         "ادو پرفیوم زنانه بالرینا مدل گود گرل Good Girl حجم 90 میلی لیتر" + strconv.Itoa(i*2),
			Slug:          "ادو-پرفیوم-زنانه-بالرینا-مدل-گود-گرل-good-girl" + strconv.Itoa(i*2),
			Sku:           "sku1000" + strconv.Itoa(i*2),
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
		})
	}

	p2 := []entities.Product{
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

		//men-pents-trousers
		{
			CategoryID:    22,
			BrandID:       5,
			Title:         "شلوار مردانه مدل بنگال کمربند دار",
			Slug:          "شلوار-مردانه-مدل-بنگال-کمربند-دار",
			Sku:           "sku2000",
			Status:        true,
			OriginalPrice: 280_000,
			SalePrice:     238_000,
			Description:   "شلوار از پارچه ی به اصطلاح بنگال تولید شده است،پارچه ی کتان بنگال پارچه ای با ظرافت بالا همراه با کشسانی نسبی مناسب می باشد که زیبایی دو چندانی در پوشیدن شلوار به شما می دهد پس اگر دنبال شلوار ضخیم میگردید ما پارچه ی بنگال را توصیه نمیکنیم.قد شلوار صد سانتی متر است،پاچه ی شلوار پاکتی است و در قسمت پاچه و کمربند مارک فلزی کار شده است،قسمت پشت کمر کش کار شده است و در جلوی کار طراحی کمربندی زیبا که شمارا از بستن کمربند بی نیاز میکند و راحتی دو چندانی را به ارمغان خواهد آورد.شلوار دارای دو جیب در بغل و یک جیب کوچک در پشت است،یک ساسون در پای چپ و یک ساسون در روی پای راست به ظاهر کلاسیکی شلوار می افزاید.رنگ شلوار مشکی است و مهمترین ویژگی آن استایل جذب و قابلیت پوشیدن با کفش کالج و تیپ رسمی و همینطور قابلیت پوشیدن با کفش اسپرت و تیپ اسپرت را دارد.",

			ProductImages: []entities.ProductImages{
				{
					Path: "2024/09/27/4444.webp",
				},
				{
					Path: "2024/09/27/5555.webp",
				},
				{
					Path: "2024/09/27/6666.webp",
				},
			},
			ProductAttributes: []entities.ProductAttribute{

				{
					//medium
					//ProductID:           4,
					AttributeID:         1,
					AttributeTitle:      "سایز",
					AttributeValueID:    2,
					AttributeValueTitle: "m",
				},
				{
					//large
					//ProductID:           4,
					AttributeID:         1,
					AttributeTitle:      "سایز",
					AttributeValueID:    3,
					AttributeValueTitle: "l",
				},
				{
					//x-large
					//ProductID:           4,
					AttributeID:         1,
					AttributeTitle:      "سایز",
					AttributeValueID:    4,
					AttributeValueTitle: "xl",
				},
			},
			ProductInventories: []entities.ProductInventory{
				{
					//ProductID: 4,
					Quantity: 25,
				},
				{
					//ProductID: 4,
					Quantity: 50,
				},
				{
					//ProductID: 4,
					Quantity: 75,
				},
			},
			ProductInventoryAttributes: []entities.ProductInventoryAttribute{
				{
					ProductID:          204,
					ProductInventoryID: 204,
					ProductAttributeID: 1, //medium
				},
				{
					ProductID:          204,
					ProductInventoryID: 205,
					ProductAttributeID: 2, //large
				},
				{
					ProductID:          204,
					ProductInventoryID: 206,
					ProductAttributeID: 3, //x-large
				},
			},
			Features: []entities.Feature{
				{
					Title: "جنس",
					Value: "بنگال",
				},
				{
					Title: "طرح",
					Value: "ساده",
				},
				{
					Title: "استایل شلوار",
					Value: "جذب",
				},
				{
					Title: "نوع فاق",
					Value: "متوسط",
				},
				{
					Title: "نحوه بسته شدن",
					Value: "یکسره",
				},
				{
					Title: "مورد استفاده",
					Value: "روزمره\n\nاداری و رسمی\n\nمهمانی\n\n",
				},
			},
		},
	}

	p = append(p, p2...)
	return p
}

func fakeAttributeAndValues() []entities.Attribute {
	return []entities.Attribute{
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
}

func fakeCategories() []entities.Category {
	id1 := uint(1)
	id2 := uint(2)

	return []entities.Category{
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
}

func fakeBrands() []entities.Brand {
	return []entities.Brand{
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

		{
			Title: "داخلی",
			Slug:  "iran",
		},
	}
}
