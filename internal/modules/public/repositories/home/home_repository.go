package home

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"math/rand"
	"shop/internal/database/mysql"
	"shop/internal/entities"
	"shop/internal/modules/public/requests"
	"shop/internal/modules/public/responses"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/helpers"
	"shop/internal/pkg/pagination"
	"shop/internal/pkg/sessions"
	"shop/internal/pkg/util"
	"strconv"
	"strings"
	"time"
)

type HomeRepository struct {
	db *gorm.DB
}

func NewHomeRepository() HomeRepository {
	return HomeRepository{
		db: mysql.Get(),
	}
}
func (h HomeRepository) GetRandomProducts(ctx context.Context, limit int) ([]entities.Product, error) {
	var products []entities.Product
	//implement Me
	return products, nil
}

func (h HomeRepository) GetLatestProducts(ctx context.Context, limit int) ([]entities.Product, error) {
	var products []entities.Product
	//todo: just load data if category.status = true and product.status=true
	err := h.db.Preload("Category").Where("status=?", true).Limit(limit).Find(&products).Error
	return products, err
}

func (h HomeRepository) GetCategories(ctx context.Context, limit int) ([]entities.Category, error) {
	var categories []entities.Category
	err := h.db.Limit(limit).Find(&categories, "status=?", true).Error

	return categories, err
}
func (h HomeRepository) GetProduct(c *gin.Context, productSku string, productSlug string) (map[string]interface{}, error) {

	type InventoryWithAttributes struct {
		InventoryID                 uint
		Quantity                    uint
		AttributeID                 uint
		AttributeTitle              string
		AttributeValueID            uint
		AttributeValueTitle         string
		ProductInventoryAttributeID uint
	}

	var product entities.Product
	aerr := h.db.WithContext(c).
		Preload("Category").
		Preload("Brand").
		Preload("ProductImages").
		Preload("Features").
		Where("sku=? and slug=? and status=true", productSku, productSlug).
		First(&product).Error

	if aerr != nil {
		return map[string]interface{}{}, aerr
	}

	var inventories []InventoryWithAttributes

	result := make(map[string]interface{})

	serr := h.db.
		WithContext(c).
		Table("product_inventories").
		Select("product_inventories.id AS inventory_id, product_inventories.quantity, product_attributes.attribute_id, attributes.title AS attribute_title, attribute_values.id AS attribute_value_id, attribute_values.value AS attribute_value_title, product_inventory_attributes.id AS product_inventory_attribute_id").
		Joins("LEFT JOIN product_inventory_attributes ON product_inventories.id = product_inventory_attributes.product_inventory_id AND product_inventory_attributes.deleted_at IS NULL").
		Joins("LEFT JOIN product_attributes ON product_inventory_attributes.product_attribute_id = product_attributes.id AND product_attributes.deleted_at IS NULL").
		Joins("LEFT JOIN attributes ON product_attributes.attribute_id = attributes.id AND attributes.deleted_at IS NULL").
		Joins("LEFT JOIN attribute_values ON product_attributes.attribute_value_id = attribute_values.id AND attribute_values.deleted_at IS NULL").
		Where("product_inventories.product_id = ? and product_inventories.deleted_at IS NULL", product.ID).
		Scan(&inventories).
		Error

	if serr != nil {
		return map[string]interface{}{}, serr
	}

	inventoryMap := make(map[uint]map[string]interface{})
	for _, inventory := range inventories {
		if _, exists := inventoryMap[inventory.InventoryID]; !exists {
			inventoryMap[inventory.InventoryID] = map[string]interface{}{
				"quantity":     inventory.Quantity,
				"inventory_id": inventory.InventoryID,
				"attributes":   []map[string]interface{}{},
			}
		}

		attributes := inventoryMap[inventory.InventoryID]["attributes"].([]map[string]interface{})
		attributes = append(attributes, map[string]interface{}{
			"attribute_id":                   inventory.AttributeID,
			"attribute_title":                inventory.AttributeTitle,
			"attribute_value_id":             inventory.AttributeValueID,
			"attribute_value_title":          inventory.AttributeValueTitle,
			"product_inventory_attribute_id": inventory.ProductInventoryAttributeID,
		})
		inventoryMap[inventory.InventoryID]["attributes"] = attributes
	}

	result["product"] = product
	result["inventories"] = inventoryMap

	return result, nil
}

func (h HomeRepository) GetProductsBy(ctx context.Context, columnName string, value any) ([]entities.Product, error) {
	var products []entities.Product
	condition := fmt.Sprintf("%s = ?", columnName)
	err := h.db.Where(condition, value).Find(&products).Error

	return products, err
}

func (h HomeRepository) GetCategoryBy(ctx context.Context, columnName string, value any) (entities.Category, error) {
	var category entities.Category
	err := h.db.Where(fmt.Sprintf("%s = ?", columnName), value).Find(&category).Error

	return category, err
}

func (h HomeRepository) NewOtp(ctx context.Context, mobile string) (entities.OTP, custom_error.CustomError) {
	var maxOTPRequestPerHour = 4
	var lastOtp entities.OTP
	var otpCount int64

	oneHourAgo := time.Now().Add(-1 * time.Hour)
	h.db.Model(entities.OTP{}).Where("mobile = ? AND created_at >= ?", mobile, oneHourAgo).Count(&otpCount)

	if otpCount >= int64(maxOTPRequestPerHour) {
		fmt.Println("---- to many request otp ---line : 77  ---- ")
		return lastOtp, custom_error.New(custom_error.OTPTooManyRequest, custom_error.OTPTooManyRequest, custom_error.OTPTooManyRequestCode)
	}

	//check under 4 min
	h.db.Where("mobile = ? AND is_expired = ? ", mobile, false).Order("created_at desc").First(&lastOtp)
	fmt.Println("-------- before check --------- : ", lastOtp)
	fmt.Println(" ******** time since ******: ", time.Since(lastOtp.CreatedAt))
	if lastOtp.ID != 0 && time.Since(lastOtp.CreatedAt) <= 4*time.Minute {
		fmt.Println("---- to soon request otp ---line : 86  ---- ")
		return lastOtp, custom_error.New(custom_error.OTPRequestTooSoon, custom_error.OTPRequestTooSoon, custom_error.OTPTooSoonCode)
	}

	newOtp := entities.OTP{
		Mobile:    mobile,
		Code:      strconv.FormatInt(util.Random4Digit(), 10),
		IsExpired: false,
	}

	if err := h.db.Create(&newOtp).Error; err != nil {
		return newOtp, custom_error.New(err.Error(), custom_error.SomethingWrongHappened, custom_error.OtpSomethingGoesWrongCode)
	}
	return newOtp, custom_error.CustomError{}
}

func (h HomeRepository) VerifyOtp(c *gin.Context, mobile string, req requests.CustomerVerifyRequest) (entities.OTP, error) {
	var otp entities.OTP
	otpCode := fmt.Sprintf("%s%s%s%s", req.N1, req.N2, req.N3, req.N4)
	fmt.Println("------ VerifyOtp : home repository : 105 : otp : ", otpCode)
	fmt.Println("------ VerifyOtp : home repository : 105 : mobile : ", mobile)
	err := h.db.Where("mobile = ? AND code = ? ", mobile, otpCode).First(&otp).Error

	fmt.Println("--- verify otp err:--- ", err)
	return otp, err
}

func (h HomeRepository) ProcessCustomerAuthenticate(c *gin.Context, mobile string) (entities.Session, error) {

	tx := h.db.Begin()
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

func (h HomeRepository) LogOut(c *gin.Context) error {
	session_id := sessions.GET(c, "session_id")

	return h.db.Where("session_id = ?", session_id).Delete(&entities.Session{}).Error
}

func (h HomeRepository) UpdateProfile(c *gin.Context, req requests.CustomerProfileRequest) error {

	var sess entities.Session
	session_id := sessions.GET(c, "session_id")
	if sessErr := h.db.Where("session_id = ? ", session_id).First(&sess).Error; sessErr != nil {
		return sessErr
	}

	var customer entities.Customer
	if cErr := h.db.First(&customer, sess.CustomerID).Error; cErr != nil {
		return cErr
	}

	if uErr := h.db.Model(&customer).
		Update("first_name", strings.TrimSpace(req.FirstName)).
		Update("last_name", strings.TrimSpace(req.LastName)).Error; uErr != nil {
		return uErr
	}

	return nil
}

func (h HomeRepository) GetMenu(ctx context.Context) ([]entities.Category, error) {
	var menu []entities.Category
	err := h.db.Preload("SubCategories", func(db *gorm.DB) *gorm.DB {
		return db.Order("priority is null ,priority ASC")
	}).Preload("SubCategories.SubCategories", func(db *gorm.DB) *gorm.DB {
		// مرتب‌سازی زیرمجموعه‌های سطح دوم
		return db.Order("priority is null ,priority ASC")
	}).Where("status = ?", true).
		Where("parent_id IS NULL").              // دریافت دسته‌های اصلی (والدین)
		Order("priority is null ,priority ASC"). // مرتب‌سازی دسته‌های اصلی
		Find(&menu).Error
	if err != nil {
		return nil, err
	}
	return menu, nil
}

func (h HomeRepository) ListProductBy(c *gin.Context, slug string) (pagination.Pagination, error) {

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
	if err := h.db.WithContext(c).Where("slug = ?", slug).First(&category).Error; err != nil {
		return pg, err
	}

	var products []entities.Product
	condition := fmt.Sprintf("category_id=%d", category.ID)

	paginateQuery, exist := pagination.Paginate(c, condition, &products, &pg, h.db)
	if !exist {
		return pg, gorm.ErrRecordNotFound
	}

	if pErr := paginateQuery(h.db).Preload("ProductImages").Where("category_id=?", category.ID).Find(&products).Error; pErr != nil {
		return pg, pErr
	}

	pg.Rows = products
	return pg, nil
}

func (h HomeRepository) InsertCart(c *gin.Context, user responses.Customer, product entities.MongoProduct, req requests.AddToCartRequest) {
	maxQuantity := uint8(2)
	//todo: set max quantity in config
	//todo:check inventories stock_reserved before insert
	//todo:check product count in the cart
	//todo: after change cart to order, we will delete cart and cartItems
	var cart entities.Cart

	//todo:check cart status
	err := h.db.
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
					ProductID:     uint(product.Product.ID),
					InventoryID:   req.InventoryID, //اگر اینونتوری صفر باشه به این معنی هست که ما برای محصول فقط موجودی ست کردیم و اون محصول دارای چند موجودی به ازای چند اتریبیوت نیست!
					Quantity:      1,
					OriginalPrice: uint(product.Product.OriginalPrice),
					SalePrice:     uint(product.Product.SalePrice),
					ProductSku:    product.Product.Sku,
					ProductTitle:  product.Product.Title,
					ProductImage:  product.Product.Images.Data[0].OriginalPath,
					ProductSlug:   product.Product.Slug,
				},
			},
		}

		h.db.Create(&cart)
		fmt.Println("~~~~~~~ [create] new cart created ,cart id is : ", cart.ID)
	} else {

		itemExist := false
		for i, cartItem := range cart.CartItems {

			if cartItem.ProductID == uint(product.Product.ID) && cartItem.InventoryID == req.InventoryID {
				if cart.CartItems[i].Quantity < maxQuantity {
					cart.CartItems[i].Quantity += 1
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
					ProductID:     uint(product.Product.ID),
					InventoryID:   req.InventoryID, //اگر اینونتوری صفر باشه به این معنی هست که ما برای محصول فقط موجودی ست کردیم و اون محصول دارای چند موجودی به ازای چند اتریبیوت نیست!
					Quantity:      1,
					OriginalPrice: uint(product.Product.OriginalPrice),
					SalePrice:     uint(product.Product.SalePrice),
					ProductSku:    product.Product.Sku,
					ProductTitle:  product.Product.Title,
					ProductImage:  product.Product.Images.Data[0].OriginalPath,
					ProductSlug:   product.Product.Slug,
				},
			}
			h.db.Model(&cart).Association("CartItems").Append(&CartItems)

		}

		h.db.Save(&cart)

	}

}

func (h HomeRepository) IncreaseCartItemCount(c *gin.Context, req requests.IncreaseCartItemQty) error {
	customer, exist := helpers.GetAuthUser(c)
	if !exist {
		return errors.New(custom_error.SomethingWrongHappened)
	}

	//return nil
	return h.db.
		Model(&entities.CartItem{}).
		Where("cart_id=?", req.CartID).
		Where("customer_id=?", customer.ID).
		Where("product_id=?", req.ProductID).
		Where("inventory_id=?", req.InventoryID).
		Where("quantity<?", 3).                                //max qty to order
		Update("quantity", gorm.Expr("quantity + ?", 1)).Error //todo:qty <3

}

func (h HomeRepository) DecreaseCartItemCount(c *gin.Context, req requests.IncreaseCartItemQty) error {
	customer, exist := helpers.GetAuthUser(c)
	if !exist {
		return errors.New(custom_error.SomethingWrongHappened)
	}

	return h.db.
		Model(&entities.CartItem{}).
		Where("cart_id=?", req.CartID).
		Where("customer_id=?", customer.ID).
		Where("product_id=?", req.ProductID).
		Where("inventory_id=?", req.InventoryID).
		Where("quantity>?", 1).
		Update("quantity", gorm.Expr("quantity - ?", 1)).
		Error

}
func (h HomeRepository) DeleteCartItem(c *gin.Context, req requests.IncreaseCartItemQty) error {
	customer, exist := helpers.GetAuthUser(c)
	if !exist {
		return errors.New(custom_error.SomethingWrongHappened)
	}

	return h.db.
		Model(&entities.CartItem{}).Unscoped().
		Where("cart_id=?", req.CartID).
		Where("customer_id=?", customer.ID).
		Where("product_id=?", req.ProductID).
		Where("inventory_id=?", req.InventoryID).
		Delete(&entities.CartItem{}).
		Error

}

func (h HomeRepository) CreateOrUpdateAddress(c *gin.Context, req requests.StoreAddressRequest) error {
	customer, ok := helpers.GetAuthUser(c)
	if !ok {
		return nil
	}

	if customer.Address.ID <= 0 {
		if err := h.db.Create(&entities.Address{
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

	if err := h.db.Model(&entities.Address{}).Where("customer_id=?", customer.ID).Updates(&entities.Address{
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

// GenerateOrderFromCart create new order and new order-item from cart and cart-item then remove cart
func (h HomeRepository) GenerateOrderFromCart(c *gin.Context) (entities.Order, error) {
	customer, ok := helpers.GetAuthUser(c)
	if !ok {
		return entities.Order{}, errors.New(custom_error.SomethingWrongHappened)
	}

	tx := h.db.WithContext(c).Begin()

	//check qty and reserve it
	for _, cartItem := range customer.Cart.CartItem.Data {
		var pInventory entities.ProductInventory

		//lock for update
		if cItemErr := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND product_id = ?", cartItem.InventoryID, cartItem.ProductID).
			First(&pInventory).
			Error; cItemErr != nil {
			tx.Rollback()
			fmt.Println("[home_repository]-[GenerateOrderFromCart]-[find-inventory]-error:", cItemErr.Error())
			return entities.Order{}, cItemErr
		}

		realQty := pInventory.Quantity - pInventory.ReservedStock
		if realQty < uint(cartItem.Quantity) {
			//out of stock
			tx.Rollback()
			return entities.Order{}, errors.New("out of stock")
		}

		//update and save reserved stock
		pInventory.ReservedStock += uint(cartItem.Quantity)
		if updateInventoryReservedStock := tx.Save(&pInventory).Error; updateInventoryReservedStock != nil {
			tx.Rollback()
			fmt.Println("[home_repository]-[GenerateOrderFromCart]-[update-reserve-stock]-error:", updateInventoryReservedStock.Error())
			return entities.Order{}, updateInventoryReservedStock
		}

	}

	//convert address struct to json
	addressJson, _ := json.Marshal(customer.Address)

	//create order
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

	if createOrderError := tx.Create(&order).Error; createOrderError != nil {
		tx.Rollback()
		fmt.Println("[home_repository]-[GenerateOrderFromCart]-[create-order]-error:", createOrderError.Error())
		return order, createOrderError
	}

	//create order-items
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
		util.PrettyJson(cartItem)
	}

	if createOrderItemsErr := tx.Create(&orderItems).Error; createOrderItemsErr != nil {
		tx.Rollback()
		fmt.Println("[home_repository]-[GenerateOrderFromCart]-[create-order-items]-error:", createOrderItemsErr.Error())
		return order, createOrderItemsErr
	}

	//Delete Cart and its CartItem
	if true {
		if deleteCartErr := h.db.WithContext(c).Unscoped().Delete(&entities.Cart{}, customer.Cart.ID).Error; deleteCartErr != nil {
			tx.Rollback()
			fmt.Println("[home_repository]-[GenerateOrderFromCart]-[delete-cart-and-cartItem]-error:", deleteCartErr.Error())
			return order, deleteCartErr
		}
	}

	tx.Commit()
	return order, nil
}

func (h HomeRepository) Release(order entities.Order, tx *gorm.DB) {
	tx.Rollback()
}

func (h HomeRepository) OrderPaidSuccessfully(c *gin.Context, order entities.Order, refID string, verified bool) (entities.Order, bool, custom_error.CustomError) {
	util.PrettyJson(order)
	log.Println("---- verify :", verified)

	//	var order entities.Order
	tx := h.db.WithContext(c).Begin()

	//if err := tx.Preload("OrderItems").Where("id=? AND amount=?", payment.OrderID).First(&order).Error; err != nil {
	//	tx.Rollback()
	//	return order, false, custom_error.New(err.Error(), custom_error.RecordNotFound, custom_error.PaymentNotFound)
	//}

	if order.OrderStatus > 0 {
		log.Println("----------order status is greater than 0 -------")
		tx.Rollback()
		log.Println("order has already been marked as paid, skipping duplicate process")
		return order, true, custom_error.New("order has already been marked as paid, skipping duplicate process", custom_error.OrderAlreadyMarkedAsPaid, custom_error.OrderMarkedAsPaid)
	}
	if verified {
		log.Println("--- x1")
		order.OrderStatus = 1 //paid successful
		if saveOrderErr := tx.Save(&order).Error; saveOrderErr != nil {
			tx.Rollback()
			return order, false, custom_error.New(saveOrderErr.Error(), custom_error.OrderChangeStatusToPaid, custom_error.OrderSavePaidStatusFailed)
		}

		order.Payment.RefID = refID
		order.Payment.Status = 1 //paid
		if updatePayment := tx.Save(&order.Payment).Error; updatePayment != nil {
			tx.Rollback()
			return order, false, custom_error.New(updatePayment.Error(), custom_error.UpdatePaymentFaileds, custom_error.UpdatePaymentFailed)
		}

	} else {
		log.Println("--- x2")
		order.OrderStatus = 2 //failed
		if saveOrderErr := tx.Save(&order).Error; saveOrderErr != nil {
			tx.Rollback()
			return order, false, custom_error.New(saveOrderErr.Error(), custom_error.UpdateOrderFaileds, custom_error.UpdateOrderFailed)
		}

		order.Payment.Status = 2 //failed
		if updatePayment := tx.Save(&order.Payment).Error; updatePayment != nil {
			tx.Rollback()
			return order, false, custom_error.New(updatePayment.Error(), custom_error.UpdatePaymentFaileds, custom_error.UpdatePaymentFailed)
		}
	}

	//decrees product inventory quantity and product inventory reserved stock
	for _, orderItem := range order.OrderItems {
		log.Println("----- x 3")
		var productInventory entities.ProductInventory
		if findProductInventoryErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("product_id=? AND id=?", orderItem.ProductID, orderItem.InventoryID).
			First(&productInventory).Error; findProductInventoryErr != nil {
			log.Println("----- x 4")

			tx.Rollback()
			return order, false, custom_error.New(findProductInventoryErr.Error(), custom_error.ProductInventoryNotFounds, custom_error.ProductInventoryNotFound)
		}

		if verified {
			log.Println("----- x 5")

			productInventory.Quantity -= orderItem.Quantity
			productInventory.ReservedStock -= orderItem.Quantity
		} else {
			log.Println("----- x 6")

			productInventory.ReservedStock -= orderItem.Quantity
		}

		//todo: update mongodb
		if updateProductInventoryErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Save(&productInventory).Error; updateProductInventoryErr != nil {
			log.Println("----- x 7")

			tx.Rollback()
			return order, false, custom_error.New(updateProductInventoryErr.Error(), custom_error.UpdateProductInventoryFaileds, custom_error.UpdateProductInventoryFailed)
		}
		log.Println("----- x 8")

	}

	tx.Commit()
	log.Println("----- x 9")

	return order, true, custom_error.CustomError{}

}

func (h HomeRepository) CreatePayment(c *gin.Context, payment entities.Payment) error {
	if err := h.db.Create(&payment).Error; err != nil {
		return err
	}
	return nil
}

func (h HomeRepository) GetPayment(c *gin.Context, authority string) (entities.Order, entities.Customer, error) {
	var payment entities.Payment
	var order entities.Order

	if err := h.db.WithContext(c).Where("authority = ?", authority).First(&payment).Error; err != nil {
		return order, entities.Customer{}, err
	}

	if orderErr := h.db.WithContext(c).Preload("OrderItems").Where("id=?", payment.OrderID).First(&order).Error; orderErr != nil {
		return order, entities.Customer{}, orderErr
	}

	var customer entities.Customer
	if customerErr := h.db.WithContext(c).Where("id=?", payment.CustomerID).First(&customer).Error; customerErr != nil {
		return order, entities.Customer{}, customerErr
	}

	order.Payment = payment
	return order, customer, nil
}

func (h HomeRepository) GetPaginatedOrders(c *gin.Context) (pagination.Pagination, error) {

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
	//if err := h.db.WithContext(c).Where("customer_id = ?", customer.ID).First(&order).Error; err != nil {
	//	return pg, err
	//}

	var orders []entities.Order
	condition := fmt.Sprintf("customer_id=%d", customer.ID)

	paginateQuery, exist := pagination.Paginate(c, condition, &orders, &pg, h.db)
	util.PrettyJson(customer)

	if !exist {
		log.Println("--------- here :", exist)
		return pg, gorm.ErrRecordNotFound
	}

	if pErr := paginateQuery(h.db).Preload("OrderItems").Where("customer_id=?", customer.ID).Find(&orders).Error; pErr != nil {
		log.Println("----------- here 2:", pErr)
		return pg, pErr
	}

	pg.Rows = orders
	return pg, nil
}

func (h HomeRepository) GetOrder(c *gin.Context, orderNumber string) (entities.Order, error) {

	customer, exists := helpers.GetAuthUser(c)
	if !exists {
		return entities.Order{}, errors.New("user must be loggedIn")
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
	if err := h.db.WithContext(c).
		Table("order_items").
		Select("product_id , inventory_id").
		Where("customer_id = ?", customer.ID).
		Scan(&productAndInventory).Error; err != nil {
		return order, nil
	}

	//حالا نتایج رو به صورت اسلایس در میاریم که مستقیم بشه درون preload استفاده کرد
	var productIDs, inventoryIDs []uint
	for _, item := range productAndInventory {
		productIDs = append(productIDs, item.ProductID)
		inventoryIDs = append(inventoryIDs, item.InventoryID)
	}

	if err := h.db.WithContext(c).
		Preload("OrderItems.Product.ProductInventoryAttributes",
			"product_inventory_attributes.product_id IN (?) AND product_inventory_attributes.product_inventory_id IN (?)",
			productIDs, inventoryIDs,
		).
		Preload("OrderItems.Product.ProductInventoryAttributes.ProductAttribute",
			"product_attributes.product_id IN (?)",
			productIDs,
		).
		Preload("Payment").
		Where("order_number=? AND customer_id = ?", orderNumber, customer.ID).
		First(&order).Error; err != nil {
		return order, err
	}

	return order, nil
}
