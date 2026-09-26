package home

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"math/rand"
	AdminUserResponse "shop/application/dto/admin"
	"shop/application/dto/web"
	"shop/bootstrap"
	"shop/domain/domain_err"
	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/infrastructure/events"
	"shop/infrastructure/repositories/product"
	"shop/interfaces/http/requests/web"
	"shop/pkg/helpers"
	"shop/pkg/pagination"
	"shop/pkg/sessions"
	"shop/pkg/util"
	"strconv"
	"strings"
	"time"
)

type HomeRepository struct {
	dep          *bootstrap.Dependencies
	eventManager *events.EventManager
}

func NewHomeRepository(dep *bootstrap.Dependencies, eventManager *events.EventManager) repositories.HomeRepositoryInterface {
	return &HomeRepository{
		dep:          dep,
		eventManager: eventManager,
	}
}

//--------------------------------------
//>>>>>>>>>>> Method >>>>>>>>>>>>>>>>>>>
//--------------------------------------

func (h *HomeRepository) GetRandomProducts(ctx context.Context, limit int) ([]*entities.Product, error) {
	//var products []*entities.Product
	//implement Me
	return nil, nil
}
func (h *HomeRepository) GetLatestProducts(ctx context.Context, limit int) ([]*entities.Product, error) {
	var products []*entities.Product
	//todo: just load data if category.status = true and product.status=published
	err := h.dep.DB.WithContext(ctx).
		Preload("Category").Where("status=?", entities.ProductStatusPublished).
		Limit(limit).Find(&products).
		Error

	return products, err
}
func (h *HomeRepository) GetCategories(ctx context.Context, limit int) ([]*entities.Category, error) {
	var categories []*entities.Category
	err := h.dep.DB.WithContext(ctx).
		Limit(limit).
		Find(&categories, "status=?", true).
		Error

	return categories, err
}

// GetProduct loads the storefront single-product payload from products.read_model (JSONB)
// plus the product's recommendations. The returned map keeps the exact shape the
// single_product.html template expects: _id / product / inventories.
func (h *HomeRepository) GetProduct(c *gin.Context, productSku string, productSlug string) (map[string]interface{}, []entities.RecommendedProduct, error) {
	var prod entities.Product
	err := h.dep.DB.WithContext(c).
		Select("id, read_model").
		Where("sku = ? AND slug = ? AND status = ?", productSku, productSlug, entities.ProductStatusPublished).
		First(&prod).
		Error
	if err != nil {
		// return the raw error so the service layer maps gorm.ErrRecordNotFound to 404
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("product-slug : ", productSlug, " | SKU :", productSku, " not found.")
		}
		return nil, nil, err
	}

	readModel, err := h.loadReadModel(c, &prod)
	if err != nil {
		return nil, nil, err
	}

	result := map[string]interface{}{
		"_id":         strconv.FormatUint(uint64(prod.ID), 10),
		"product":     readModel.Product,
		"inventories": readModel.Inventories,
	}

	recommendations := h.loadRecommendations(c, prod.ID)
	return result, recommendations, nil
}

// loadReadModel returns the stored JSONB read model; for legacy rows without one
// it rebuilds (and persists) the model on the fly.
func (h *HomeRepository) loadReadModel(c *gin.Context, prod *entities.Product) (*entities.ProductReadModel, error) {
	if len(prod.ReadModel) > 0 && string(prod.ReadModel) != "{}" {
		var rm entities.ProductReadModel
		if err := json.Unmarshal(prod.ReadModel, &rm); err == nil && rm.Product.ID != 0 {
			return &rm, nil
		}
	}

	// legacy row: build the read model now and persist it
	if err := product.SyncReadModel(c, h.dep.DB, prod.ID); err != nil {
		return nil, err
	}
	var refreshed entities.Product
	if err := h.dep.DB.WithContext(c).Select("id, read_model").First(&refreshed, prod.ID).Error; err != nil {
		return nil, err
	}
	var rm entities.ProductReadModel
	if err := json.Unmarshal(refreshed.ReadModel, &rm); err != nil {
		return nil, err
	}
	return &rm, nil
}

// loadRecommendations reads product_recommendations and builds the storefront
// recommendation list (single_product.html: $rec.Product.*).
func (h *HomeRepository) loadRecommendations(c *gin.Context, productID uint) []entities.RecommendedProduct {
	var recs []entities.ProductRecommendation
	if err := h.dep.DB.WithContext(c).
		Where("product_id = ?", productID).
		Find(&recs).
		Error; err != nil || len(recs) == 0 {
		return nil
	}

	ids := make([]uint, 0, len(recs))
	for _, r := range recs {
		ids = append(ids, r.RecommendedProductID)
	}

	var products []entities.Product
	if err := h.dep.DB.WithContext(c).
		Select("id, read_model").
		Where("id IN ?", ids).
		Find(&products).
		Error; err != nil {
		return nil
	}

	out := make([]entities.RecommendedProduct, 0, len(products))
	for i := range products {
		var rm entities.ProductReadModel
		if len(products[i].ReadModel) == 0 {
			continue
		}
		if err := json.Unmarshal(products[i].ReadModel, &rm); err != nil || rm.Product.ID == 0 {
			continue
		}
		out = append(out, entities.RecommendedProduct{Product: rm.Product})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// GetProductByID loads a product (with images) for cart operations.
func (h *HomeRepository) GetProductByID(c *gin.Context, productID uint) (*entities.Product, error) {
	var prod entities.Product
	err := h.dep.DB.WithContext(c).
		Preload("ProductImages").
		First(&prod, productID).
		Error
	return &prod, err
}

func (h *HomeRepository) GetProductsBy(ctx context.Context, columnName string, value any) ([]*entities.Product, error) {
	var products []*entities.Product
	condition := fmt.Sprintf("%s = ?", columnName)
	err := h.dep.DB.WithContext(ctx).
		Where(condition, value).
		Find(&products).
		Error

	return products, err
}
func (h *HomeRepository) GetCategoryBy(ctx context.Context, columnName string, value any) (*entities.Category, error) {
	var category entities.Category
	err := h.dep.DB.WithContext(ctx).
		Where(fmt.Sprintf("%s = ?", columnName), value).
		Find(&category).
		Error

	return &category, err
}
func (h *HomeRepository) NewOtp(ctx context.Context, mobile string) (*entities.OTP, domain_err.CustomError) {
	var maxOTPRequestPerHour = 4
	var lastOtp entities.OTP
	var otpCount int64

	oneHourAgo := time.Now().Add(-1 * time.Hour)
	h.dep.DB.Model(entities.OTP{}).WithContext(ctx).Where("mobile = ? AND created_at >= ?", mobile, oneHourAgo).Count(&otpCount)

	if otpCount >= int64(maxOTPRequestPerHour) {
		fmt.Println("---- to many request otp ---line : 77  ---- ")
		return nil, domain_err.New(domain_err.OTPTooManyRequest, domain_err.OTPTooManyRequest, domain_err.OTPTooManyRequestCode)
	}

	//check under 4 min
	h.dep.DB.WithContext(ctx).Where("mobile = ? AND is_expired = ? ", mobile, false).Order("created_at desc").First(&lastOtp)
	fmt.Println("-------- before check --------- : ", lastOtp)
	fmt.Println(" ******** time since ******: ", time.Since(lastOtp.CreatedAt))
	if lastOtp.ID != 0 && time.Since(lastOtp.CreatedAt) <= 4*time.Minute {
		fmt.Println("---- to soon request otp ---line : 86  ---- ")
		return nil, domain_err.New(domain_err.OTPRequestTooSoon, domain_err.OTPRequestTooSoon, domain_err.OTPTooSoonCode)
	}

	newOtp := entities.OTP{
		Mobile:    mobile,
		Code:      strconv.FormatInt(util.Random4Digit(), 10),
		IsExpired: false,
	}

	if err := h.dep.DB.WithContext(ctx).Create(&newOtp).Error; err != nil {
		return nil, domain_err.New(err.Error(), domain_err.SomethingWrongHappened, domain_err.OtpSomethingGoesWrongCode)
	}
	return &newOtp, domain_err.CustomError{}
}
func (h *HomeRepository) VerifyOtp(c *gin.Context, mobile string, req *requests.CustomerVerifyRequest) (*entities.OTP, error) {
	var otp entities.OTP
	otpCode := fmt.Sprintf("%s%s%s%s", req.N1, req.N2, req.N3, req.N4)
	fmt.Println("------ VerifyOtp : home repository : 105 : otp : ", otpCode)
	fmt.Println("------ VerifyOtp : home repository : 105 : mobile : ", mobile)
	err := h.dep.DB.WithContext(c).Where("mobile = ? AND code = ? ", mobile, otpCode).First(&otp).Error

	fmt.Println("--- verify otp err:--- ", err)
	return &otp, err
}
func (h *HomeRepository) ProcessCustomerAuthenticate(c *gin.Context, mobile string) (entities.Session, error) {

	tx := h.dep.DB.Begin()
	if tx.Error != nil {
		return entities.Session{}, tx.Error
	}

	//find customer
	var customer entities.Customer
	customerErr := tx.WithContext(c).Where("mobile = ? ", mobile).First(&customer).Error
	if customerErr != nil {
		if errors.Is(customerErr, gorm.ErrRecordNotFound) {

			//customer not found , fill customer
			customer = entities.Customer{
				Mobile:    mobile,
				FirstName: "",
				LastName:  "",
				Active:    true,
			}

			//store customer in db
			if createCustomerErr := tx.WithContext(c).Create(&customer).Error; createCustomerErr != nil {
				tx.Rollback()
				fmt.Println("create new customer err : ", createCustomerErr.Error())
				return entities.Session{}, createCustomerErr
			}
		} else {
			//some internal error
			tx.Rollback()
			fmt.Println("--- database internal err : ", customerErr.Error())
			return entities.Session{}, customerErr
		}
	}

	//generate uuid
	uuidValue, uuidErr := uuid.NewUUID()
	if uuidErr != nil {
		tx.Rollback()
		fmt.Println("---- generate uuid was failed :", uuidErr)
		return entities.Session{}, uuidErr
	}

	//fill session
	sess := entities.Session{
		Mobile:     customer.Mobile,
		CustomerID: customer.ID,
		SessionID:  uuidValue.String(),
		IsActive:   true,
		ExpiredAt:  time.Now().Add(365 * (24 * time.Hour)),
	}

	//store session in db
	if sessCreateErr := tx.WithContext(c).Create(&sess).Error; sessCreateErr != nil {
		fmt.Println("---- create a session failed : ", sessCreateErr)
		tx.Rollback()
		return entities.Session{}, sessCreateErr
	}

	//commit tx
	if commitErr := tx.Commit().Error; commitErr != nil {
		fmt.Println("--- commit was failed : ", commitErr)
		tx.Rollback()
		return entities.Session{}, commitErr
	}
	return sess, nil

}
func (h *HomeRepository) LogOut(c *gin.Context) error {
	sessionId := sessions.GET(c, "session_id")

	return h.dep.DB.Where("session_id = ?", sessionId).
		Delete(&entities.Session{}).
		Error
}
func (h *HomeRepository) UpdateProfile(c *gin.Context, req *requests.CustomerProfileRequest) error {

	var sess entities.Session
	sessionId := sessions.GET(c, "session_id")
	if sessErr := h.dep.DB.Where("session_id = ? ", sessionId).First(&sess).Error; sessErr != nil {
		return sessErr
	}

	var customer entities.Customer
	if cErr := h.dep.DB.First(&customer, sess.CustomerID).Error; cErr != nil {
		return cErr
	}

	if uErr := h.dep.DB.Model(&customer).
		Update("first_name", strings.TrimSpace(req.FirstName)).
		Update("last_name", strings.TrimSpace(req.LastName)).Error; uErr != nil {
		return uErr
	}

	return nil
}
func (h *HomeRepository) GetMenu(ctx context.Context) ([]*entities.Category, error) {
	var menu []*entities.Category
	err := h.dep.DB.WithContext(ctx).
		Preload("SubCategories", func(db *gorm.DB) *gorm.DB {
			return db.Order("priority is null ,priority ASC")
		}).
		Preload("SubCategories.SubCategories", func(db *gorm.DB) *gorm.DB {
			// مرتب‌سازی زیرمجموعه‌های سطح دوم
			return db.Order("priority is null ,priority ASC")
		}).
		Where("status = ?", true).
		Where("parent_id IS NULL").              // دریافت دسته‌های اصلی (والدین)
		Order("priority is null ,priority ASC"). // مرتب‌سازی دسته‌های اصلی
		Find(&menu).
		Error

	if err != nil {
		return nil, err
	}
	return menu, nil
}
func (h *HomeRepository) ListProductBy(c *gin.Context, slug string) (pagination.Pagination, error) {

	// Convert query parameters from string to int
	limitStr := c.Query("limit")
	pageStr := c.Query("page")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 { // Default to 10 if invalid
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 { // Default to 1 if invalid
		page = 1
	}
	var pg = pagination.Pagination{
		Limit: limit,
		Page:  page,
	}

	var category entities.Category
	if err := h.dep.DB.WithContext(c).Where("slug = ?", slug).First(&category).Error; err != nil {
		return pg, err
	}

	var products []*entities.Product
	condition := fmt.Sprintf("category_id=%d", category.ID)

	paginateQuery, exist := pagination.Paginate(c, condition, &products, &pg, h.dep.DB)
	if !exist {
		return pg, gorm.ErrRecordNotFound
	}

	if pErr := paginateQuery(h.dep.DB).
		Preload("Category").
		Preload("ProductImages").
		Where("category_id=?", category.ID).
		Find(&products).Error; pErr != nil {
		return pg, pErr
	}

	pg.Rows = AdminUserResponse.ToProducts(products)

	return pg, nil
}
func (h *HomeRepository) InsertCart(c *gin.Context, user responses.Customer, product *entities.Product, req requests.AddToCartRequest) {
	maxQuantity := uint8(2)
	//todo: set max quantity in config
	//todo:check inventories stock_reserved before insert
	//todo:check product count in the cart
	//todo: after change cart to order, we will delete cart and cartItems
	var cart entities.Cart

	imagePath := ""
	if len(product.ProductImages) > 0 {
		imagePath = product.ProductImages[0].Path
	}

	// prices live on the selected variant (NULL columns inherit the product)
	originalPrice, salePrice := h.cartPrices(c, product, req.InventoryID)

	//todo:check cart status
	err := h.dep.DB.
		WithContext(c).
		Preload("CartItems").
		Where("customer_id = ? AND status = ? ", user.ID, 0).
		First(&cart).
		Error

	if err != nil {
		//not found ,so we will create it
		cart := entities.Cart{
			CustomerID: user.ID,
			Status:     0,
			CartItems: []entities.CartItem{
				{
					CustomerID:    user.ID,
					ProductID:     product.ID,
					InventoryID:   req.InventoryID, //اگر اینونتوری صفر باشه به این معنی هست که ما برای محصول فقط موجودی ست کردیم و اون محصول دارای چند موجودی به ازای چند اتریبیوت نیست!
					Quantity:      1,
					OriginalPrice: originalPrice,
					SalePrice:     salePrice,
					ProductSku:    product.Sku,
					ProductTitle:  product.Title,
					ProductImage:  imagePath,
					ProductSlug:   product.Slug,
				},
			},
		}

		h.dep.DB.Create(&cart)
		fmt.Println("~~~~~~~ [create] new cart created ,cart id is : ", cart.ID)
	} else {

		itemExist := false
		for i, cartItem := range cart.CartItems {

			if cartItem.ProductID == product.ID && cartItem.InventoryID == req.InventoryID {
				if cart.CartItems[i].Quantity < maxQuantity {
					cart.CartItems[i].Quantity += 1
					// persist the incremented child explicitly: Save(&cart) does not
					// cascade updates to loaded has-many associations
					if uErr := h.dep.DB.Save(&cart.CartItems[i]).Error; uErr != nil {
						fmt.Println("[InsertCart]-[increment]-err:", uErr)
						return
					}
				} else {
					return
				}

				itemExist = true
				break
			}
		}

		if !itemExist {
			//cartItem was not exist
			CartItems := []entities.CartItem{
				{
					CustomerID:    user.ID,
					ProductID:     product.ID,
					InventoryID:   req.InventoryID, //اگر اینونتوری صفر باشه به این معنی هست که ما برای محصول فقط موجودی ست کردیم و اون محصول دارای چند موجودی به ازای چند اتریبیوت نیست!
					Quantity:      1,
					OriginalPrice: originalPrice,
					SalePrice:     salePrice,
					ProductSku:    product.Sku,
					ProductTitle:  product.Title,
					ProductImage:  imagePath,
					ProductSlug:   product.Slug,
				},
			}
			h.dep.DB.Model(&cart).Association("CartItems").Append(&CartItems)

		}

		h.dep.DB.Save(&cart)

	}
}

func (h *HomeRepository) IncreaseCartItemCount(c *gin.Context, req *requests.IncreaseCartItemQty) error {
	log.Printf("data : %+v \n", req)

	customer, exist := helpers.GetAuthUser(c)
	if !exist {
		fmt.Println("----1")
		return errors.New(domain_err.SomethingWrongHappened)
	}
	fmt.Println("----2")

	var currentQty int
	checkCartQtyErr := h.dep.DB.
		WithContext(c).
		Model(&entities.CartItem{}).
		Select("SUM(quantity)").
		Where("cart_id = ?", req.CartID).
		Where("customer_id = ?", customer.ID).
		Where("product_id = ?", req.ProductID).
		Where("inventory_id = ?", req.InventoryID).
		Scan(&currentQty).Error

	if checkCartQtyErr != nil {
		fmt.Println("----3 :current qty:", currentQty)
		return errors.New(domain_err.SomethingWrongHappened)
	}

	if currentQty >= 3 {
		fmt.Println("----4")
		return domain_err.QuantityExceedsLimit
	}

	//check inventory
	var productInventory entities.ProductVariant
	err := h.dep.DB.WithContext(c).
		Where("id = ? AND product_id = ?", req.InventoryID, req.ProductID).
		First(&productInventory).Error

	if err != nil {
		fmt.Println("----5")
		return errors.New(domain_err.SomethingWrongHappened)
	}
	realQty := productInventory.Stock - productInventory.ReservedStock
	fmt.Println("----6 : real qty:", realQty)

	// 2<3 || 3<3+1
	if realQty < uint(currentQty) || realQty < uint(currentQty)+1 {
		fmt.Println("----7")
		//out of stock
		return domain_err.OutOfStock
	}

	fmt.Println("----8")

	//return nil
	qtyExceedLimitErr := h.dep.DB.
		Model(&entities.CartItem{}).
		Where("cart_id=?", req.CartID).
		Where("customer_id=?", customer.ID).
		Where("product_id=?", req.ProductID).
		Where("inventory_id=?", req.InventoryID).
		Where("quantity<?", 3).                                //max qty to order
		Update("quantity", gorm.Expr("quantity + ?", 1)).Error //todo:qty <3
	if qtyExceedLimitErr != nil {
		return domain_err.QuantityExceedsLimit
	}

	return nil
}
func (h *HomeRepository) DecreaseCartItemCount(c *gin.Context, req *requests.IncreaseCartItemQty) error {
	customer, exist := helpers.GetAuthUser(c)
	if !exist {
		return errors.New(domain_err.SomethingWrongHappened)
	}

	return h.dep.DB.
		Model(&entities.CartItem{}).
		Where("cart_id=?", req.CartID).
		Where("customer_id=?", customer.ID).
		Where("product_id=?", req.ProductID).
		Where("inventory_id=?", req.InventoryID).
		Where("quantity>?", 1).
		Update("quantity", gorm.Expr("quantity - ?", 1)).
		Error

}
func (h *HomeRepository) DeleteCartItem(c *gin.Context, req *requests.IncreaseCartItemQty) error {
	customer, exist := helpers.GetAuthUser(c)
	if !exist {
		return errors.New(domain_err.SomethingWrongHappened)
	}

	return h.dep.DB.
		Model(&entities.CartItem{}).Unscoped().
		Where("cart_id=?", req.CartID).
		Where("customer_id=?", customer.ID).
		Where("product_id=?", req.ProductID).
		Where("inventory_id=?", req.InventoryID).
		Delete(&entities.CartItem{}).
		Error

}
func (h *HomeRepository) CreateOrUpdateAddress(c *gin.Context, req *requests.StoreAddressRequest) error {
	customer, ok := helpers.GetAuthUser(c)
	if !ok {
		return nil
	}

	if customer.Address.ID <= 0 {
		if err := h.dep.DB.Create(&entities.Address{
			CustomerID:         customer.ID,
			ReceiverName:       req.ReceiverName,
			ReceiverMobile:     req.ReceiverMobile,
			ReceiverAddress:    req.ReceiverAddress,
			ReceiverPostalCode: req.ReceiverPostalCode,
		}).Error; err != nil {
			util.PrettyJson(err)
			return errors.New("خطا در ذخیره آدرس")
		}
	}

	if err := h.dep.DB.Model(&entities.Address{}).
		Where("customer_id=?", customer.ID).
		Updates(&entities.Address{
			CustomerID:         customer.ID,
			ReceiverName:       req.ReceiverName,
			ReceiverMobile:     req.ReceiverMobile,
			ReceiverAddress:    req.ReceiverAddress,
			ReceiverPostalCode: req.ReceiverPostalCode,
		}).Error; err != nil {
		util.PrettyJson(err)
		return errors.New("خطا در بروزرسانی آدرس")
	}

	return nil
}

func releaseLocks(ctx context.Context, redisClient *redis.Client, lockKeys []string) {
	if len(lockKeys) > 0 {
		redisClient.Del(ctx, lockKeys...)
	}
}

// retryWithBackoff attempts an operation with retries and exponential backoff
func retryWithBackoff(attempts int, delay time.Duration, operation func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = operation(); err == nil {
			return nil
		}
		time.Sleep(delay)
		delay *= 2 // Exponential backoff
	}
	return err
}

// GenerateOrderFromCart create new order and new order-item from cart and cart-item then remove cart
// cartPrices resolves the cart item prices for the selected variant:
// original = the crossed-out list price, sale = the price the customer pays
// (discount_price wins when it is a real discount — golden rule #3).
// NULL variant price columns inherit the parent product price.
func (h *HomeRepository) cartPrices(c *gin.Context, product *entities.Product, inventoryID uint) (uint, uint) {
	originalPrice := product.OriginalPrice
	salePrice := product.SalePrice

	if inventoryID == 0 {
		return originalPrice, salePrice
	}

	var variant entities.ProductVariant
	if err := h.dep.DB.WithContext(c).
		Where("id = ? AND product_id = ?", inventoryID, product.ID).
		First(&variant).Error; err != nil {
		return originalPrice, salePrice
	}

	if variant.Price != nil {
		originalPrice = *variant.Price
	}
	if variant.SalePrice != nil {
		salePrice = *variant.SalePrice
	}
	if variant.DiscountPrice != nil && *variant.DiscountPrice > 0 && *variant.DiscountPrice < salePrice {
		salePrice = *variant.DiscountPrice
	}
	return originalPrice, salePrice
}

func (h *HomeRepository) GenerateOrderFromCart(c *gin.Context) (orderModel *entities.Order, inventoryID uint, GenerateOrderErr error) {
	customer, ok := helpers.GetAuthUser(c)
	if !ok {
		return nil, inventoryID, errors.New(domain_err.SomethingWrongHappened)
	}

	// start transaction
	tx := h.dep.DB.WithContext(c).Begin()

	// store redis lock keys
	lockKeys := make([]string, 0)

	//check qty and reserve it
	touchedProducts := make(map[uint]struct{})
	for _, cartItem := range customer.Cart.CartItem.Data {

		// generate keys to store in redis -> e.g. key "lock:inventory:203"
		lockKey := fmt.Sprintf("lock:inventory:%d", cartItem.InventoryID)

		// store `lockKey` in redis
		lockErr := retryWithBackoff(3, 100*time.Millisecond,
			func() error {
				locked, redisErr := h.dep.RedisClient.SetNX(c, lockKey, "locked", 5*time.Second).Result()
				if redisErr != nil {
					return redisErr
				}
				if !locked {
					return domain_err.InventoryLockedByAnotherOne
				}
				lockKeys = append(lockKeys, lockKey)
				return nil
			})

		if lockErr != nil {
			releaseLocks(c, h.dep.RedisClient, lockKeys)
			return nil, inventoryID, lockErr
		}

		var pInventory entities.ProductVariant

		// find specific inventory
		findErr := retryWithBackoff(3, 100*time.Millisecond,
			func() error {
				return tx.WithContext(c).
					Where("id = ? AND product_id = ?", cartItem.InventoryID, cartItem.ProductID).
					First(&pInventory).
					Error
			})
		// not found
		if findErr != nil {
			tx.Rollback()
			releaseLocks(c, h.dep.RedisClient, lockKeys)

			return nil, inventoryID, findErr
		}

		// real inventory quantity
		realQty := pInventory.Stock - pInventory.ReservedStock

		// out of stock
		if realQty < uint(cartItem.Quantity) {
			tx.Rollback()
			releaseLocks(c, h.dep.RedisClient, lockKeys)

			return nil, pInventory.ID, domain_err.OutOfStock
		}

		//var finalQty uint
		pInventory.ReservedStock += uint(cartItem.Quantity)

		//update and save reserved stock
		updateInventoryReservedStock := retryWithBackoff(3, 100*time.Millisecond,
			func() error {
				return tx.Save(&pInventory).Error
			})

		//update and save reserved stock (rollback)
		if updateInventoryReservedStock != nil {
			tx.Rollback()
			releaseLocks(c, h.dep.RedisClient, lockKeys)

			return nil, pInventory.ID, updateInventoryReservedStock
		}

		touchedProducts[pInventory.ProductID] = struct{}{}
	}

	//convert address struct to json to store in order
	addressJson, _ := json.Marshal(customer.Address)

	// prepare order entity
	order := entities.Order{
		CustomerID:         customer.ID,
		OrderNumber:        strconv.Itoa(rand.Intn(9999999)),
		PaymentStatus:      entities.PaymentPending, //pending
		TotalOriginalPrice: customer.Cart.CartItem.TotalOriginalPrice,
		TotalSalePrice:     customer.Cart.CartItem.TotalSalePrice,
		Discount:           0,
		OrderStatus:        entities.OrderPending, //pending
		Address:            string(addressJson),
	}

	// store order in db
	createOrderError := retryWithBackoff(3, 100*time.Millisecond,
		func() error {
			return tx.Create(&order).Error
		})

	// store order failed
	if createOrderError != nil {
		tx.Rollback()
		releaseLocks(c, h.dep.RedisClient, lockKeys)

		return nil, inventoryID, createOrderError
	}

	// prepare order-item (rollback)
	var orderItems []entities.OrderItem
	for _, cartItem := range customer.Cart.CartItem.Data {
		orderItems = append(orderItems, entities.OrderItem{
			CustomerID:         customer.ID,
			OrderID:            order.ID,
			ProductID:          cartItem.ProductID,
			InventoryID:        cartItem.InventoryID,
			Quantity:           uint(cartItem.Quantity),
			OriginalPrice:      cartItem.OriginalPrice,
			SalePrice:          cartItem.SalePrice,
			TotalOriginalPrice: cartItem.OriginalPrice * uint(cartItem.Quantity),
			TotalSalePrice:     cartItem.SalePrice * uint(cartItem.Quantity),
		})

	}

	// store order-item in db
	createOrderItemsErr := retryWithBackoff(3, 100*time.Millisecond,
		func() error {
			return tx.Create(&orderItems).Error
		})

	// store order-item failed(rollback)
	if createOrderItemsErr != nil {
		tx.Rollback()
		releaseLocks(c, h.dep.RedisClient, lockKeys)

		fmt.Println("[home_repository]-[GenerateOrderFromCart]-[create-order-items]-error:", createOrderItemsErr.Error())
		return nil, inventoryID, createOrderItemsErr
	}

	//Delete Cart and its CartItem
	if true {
		deleteCartErr := retryWithBackoff(3, 100*time.Millisecond,
			func() error {
				return h.dep.DB.WithContext(c).Unscoped().Delete(&entities.Cart{}, customer.Cart.ID).Error
			},
		)

		if deleteCartErr != nil {
			tx.Rollback()
			releaseLocks(c, h.dep.RedisClient, lockKeys)

			return nil, inventoryID, deleteCartErr
		}
	}

	// release redis locks
	defer releaseLocks(c, h.dep.RedisClient, lockKeys)

	tx.Commit()

	//reserved stock changed — refresh pricing aggregates after commit
	for pid := range touchedProducts {
		if pricingErr := product.RefreshProductAggregates(c, h.dep.DB, pid); pricingErr != nil {
			util.Trace(pricingErr)
		}
	}

	return &order, inventoryID, nil
}

func (h *HomeRepository) OrderPaidSuccessfully(c *gin.Context, order *entities.Order, refID string, verified bool) (*entities.Order, bool, domain_err.CustomError) {

	tx := h.dep.DB.WithContext(c).Begin()

	//if err := tx.Preload("OrderItems").Where("id=? AND amount=?", payment.OrderID).First(&order).Error; err != nil {
	//	tx.Rollback()
	//	return order, false, domain_err.New(err.Error(), domain_err.RecordNotFound, domain_err.PaymentNotFound)
	//}

	if order.OrderStatus > 0 {
		log.Println("----------order status is greater than 0 -------")
		tx.Rollback()
		log.Println("order has already been marked as paid, skipping duplicate process")
		return order, true, domain_err.New("order has already been marked as paid, skipping duplicate process", domain_err.OrderAlreadyMarkedAsPaid, domain_err.OrderMarkedAsPaid)
	}
	if verified {
		log.Println("--- x1")

		order.OrderStatus = entities.OrderConfirmed //paid successful
		order.PaymentStatus = int(entities.OrderConfirmed)
		if saveOrderErr := tx.Save(&order).Error; saveOrderErr != nil {
			tx.Rollback()
			return order, false, domain_err.New(saveOrderErr.Error(), domain_err.OrderChangeStatusToPaid, domain_err.OrderSavePaidStatusFailed)
		}

		order.Payment.RefID = refID
		order.Payment.Status = int(entities.OrderConfirmed) //paid successful
		if updatePayment := tx.Save(&order.Payment).Error; updatePayment != nil {
			tx.Rollback()
			return order, false, domain_err.New(updatePayment.Error(), domain_err.UpdatePaymentFaileds, domain_err.UpdatePaymentFailed)
		}

	} else {
		log.Println("--- x2")
		order.OrderStatus = entities.OrderCancelled        //لغو شده
		order.PaymentStatus = int(entities.OrderCancelled) //لغو شده
		if saveOrderErr := tx.Save(&order).Error; saveOrderErr != nil {
			tx.Rollback()
			return order, false, domain_err.New(saveOrderErr.Error(), domain_err.UpdateOrderFaileds, domain_err.UpdateOrderFailed)
		}

		order.Payment.Status = int(entities.OrderCancelled) //لغو شده
		if updatePayment := tx.Save(&order.Payment).Error; updatePayment != nil {
			tx.Rollback()
			return order, false, domain_err.New(updatePayment.Error(), domain_err.UpdatePaymentFaileds, domain_err.UpdatePaymentFailed)
		}
	}

	//decrees product inventory quantity and product inventory reserved stock
	lockKeys := make([]string, 0)
	touchedProducts := make(map[uint]struct{})

	for _, orderItem := range order.OrderItems {
		lockKey := fmt.Sprintf("lock:inventory:%d", orderItem.InventoryID)
		lockErr := retryWithBackoff(3, 100*time.Millisecond, func() error {
			locked, redisLockErr := h.dep.RedisClient.SetNX(c, lockKey, "locked", 1*time.Second).Result()
			if redisLockErr != nil {
				return redisLockErr
			}
			if !locked {
				return domain_err.InventoryLockedByAnotherOne
			}
			lockKeys = append(lockKeys, lockKey)
			return nil
		})
		if lockErr != nil {
			fmt.Println("failed to to lock keys in OrderPaidSuccessfully")
			releaseLocks(c, h.dep.RedisClient, lockKeys)
			return nil, false, domain_err.CustomError{}
		}

		log.Println("----- x 3")
		var productInventory entities.ProductVariant
		findProductInventoryErr :=
			retryWithBackoff(3, 100*time.Millisecond,
				func() error {
					log.Println("----- x 4")
					return tx.
						//Clauses(clause.Locking{Strength: "UPDATE"}).
						WithContext(c).
						Where("product_id=? AND id=?", orderItem.ProductID, orderItem.InventoryID).
						First(&productInventory).Error
				})
		if findProductInventoryErr != nil {
			tx.Rollback()
			releaseLocks(c, h.dep.RedisClient, lockKeys)
			return order, false, domain_err.New(findProductInventoryErr.Error(), domain_err.ProductInventoryNotFounds, domain_err.ProductInventoryNotFound)
		}

		if verified {
			log.Println("----- x 5")
			productInventory.Stock -= orderItem.Quantity
			productInventory.ReservedStock -= orderItem.Quantity
		} else {
			log.Println("----- x 6")
			productInventory.ReservedStock -= orderItem.Quantity
		}

		// update Product Inventory
		updateProductInventoryErr :=
			retryWithBackoff(3, 100*time.Millisecond, func() error {
				log.Println("----- x 7")
				return tx.Save(&productInventory).Error
			})

		// update product Inventory
		if updateProductInventoryErr != nil {
			log.Println("----- x 7-1")
			tx.Rollback()
			releaseLocks(c, h.dep.RedisClient, lockKeys)
			return order, false, domain_err.New(updateProductInventoryErr.Error(), domain_err.UpdateProductInventoryFaileds, domain_err.UpdateProductInventoryFailed)
		} else {

			// stock changed for this product — refresh happens after commit
			touchedProducts[orderItem.ProductID] = struct{}{}
		}
		log.Println("----- x 8")
	}
	defer releaseLocks(c, h.dep.RedisClient, lockKeys)
	tx.Commit()
	log.Println("----- x 9")

	//refresh pricing aggregates then the read model once the stock change is
	//committed (the Typesense document mirrors the aggregates)
	for pid := range touchedProducts {
		if pricingErr := product.RefreshProductAggregates(c, h.dep.DB, pid); pricingErr != nil {
			util.Trace(pricingErr)
		}
		if syncErr := product.SyncReadModel(c, h.dep.DB, pid); syncErr != nil {
			util.Trace(syncErr)
		}
	}

	return order, true, domain_err.CustomError{}

}

func (h *HomeRepository) CreatePayment(c *gin.Context, payment *entities.Payment) error {

	err := retryWithBackoff(3, 100*time.Millisecond,
		func() error {
			return h.dep.DB.WithContext(c).Create(payment).Error
		})

	if err != nil {
		fmt.Println("error while creating payment :", err)
		return domain_err.InternalServerErr
	}
	return nil
}

func (h *HomeRepository) GetPayment(c *gin.Context, authority string) (*entities.Order, entities.Customer, error) {
	var payment entities.Payment
	var order entities.Order

	if err := h.dep.DB.WithContext(c).Where("authority = ?", authority).First(&payment).Error; err != nil {
		return nil, entities.Customer{}, err
	}

	if orderErr := h.dep.DB.WithContext(c).Preload("OrderItems").Where("id=?", payment.OrderID).First(&order).Error; orderErr != nil {
		return nil, entities.Customer{}, orderErr
	}

	var customer entities.Customer
	if customerErr := h.dep.DB.WithContext(c).Where("id=?", payment.CustomerID).First(&customer).Error; customerErr != nil {
		return nil, entities.Customer{}, customerErr
	}

	order.Payment = &payment
	return &order, customer, nil
}
func (h *HomeRepository) GetPaginatedOrders(c *gin.Context) (pagination.Pagination, error) {

	customer, exists := helpers.GetAuthUser(c)
	if !exists {
		return pagination.Pagination{}, errors.New("user must be logged In")
	}

	// Convert query parameters from string to int
	limitStr := c.Query("limit")
	pageStr := c.Query("page")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 { // Default to 10 if invalid
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 { // Default to 1 if invalid
		page = 1
	}
	var pg = pagination.Pagination{
		Limit: limit,
		Page:  page,
	}

	//var order entities.Order
	//if err := h.dep.DB.WithContext(c).Where("customer_id = ?", customer.ID).First(&order).Error; err != nil {
	//	return pg, err
	//}

	var orders []*entities.Order
	condition := fmt.Sprintf("customer_id=%d", customer.ID)

	paginateQuery, exist := pagination.Paginate(c, condition, &orders, &pg, h.dep.DB)

	if !exist {
		return pg, gorm.ErrRecordNotFound
	}

	if pErr := paginateQuery(h.dep.DB).Preload("OrderItems").Where("customer_id=?", customer.ID).Find(&orders).Error; pErr != nil {
		return pg, pErr
	}

	pg.Rows = AdminUserResponse.ToOrders(orders)

	return pg, nil
}
func (h *HomeRepository) GetOrder(c *gin.Context, orderNumber string) (*entities.Order, error) {

	customer, exists := helpers.GetAuthUser(c)
	if !exists {
		return nil, errors.New("user must be loggedIn")
	}

	//--
	var order entities.Order

	//برای گرفتن دیتای جدول
	//product_attributes
	//مجبور هستیم که ابتدا دو ستون
	//product_id , inventory_id
	//که در جدول order_items هستند
	// رو بدست بیاریم و بعد نتایج اونها رو درون preload استفاده کنیم

	var productAndInventory []struct {
		ProductID   uint
		InventoryID uint
	}
	if err := h.dep.DB.WithContext(c).
		Table("order_items").
		Select("product_id , inventory_id").
		Where("customer_id = ?", customer.ID).
		Scan(&productAndInventory).Error; err != nil {
		return nil, err
	}

	//حالا نتایج رو به صورت اسلایس در میاریم که مستقیم بشه درون preload استفاده کرد
	var productIDs, inventoryIDs []uint
	for _, item := range productAndInventory {
		productIDs = append(productIDs, item.ProductID)
		inventoryIDs = append(inventoryIDs, item.InventoryID)
	}

	if err := h.dep.DB.WithContext(c).
		Preload("OrderItems.Product.VariantAttributeValues",
			"variant_attribute_values.product_id IN (?) AND variant_attribute_values.variant_id IN (?)",
			productIDs, inventoryIDs,
		).
		Preload("OrderItems.Product.VariantAttributeValues.AttributeValue").
		Preload("Payment").
		Where("order_number=? AND customer_id = ?", orderNumber, customer.ID).
		First(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}
