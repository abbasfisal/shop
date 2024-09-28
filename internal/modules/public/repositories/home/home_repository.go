package home

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"shop/internal/database/mysql"
	"shop/internal/entities"
	"shop/internal/modules/public/requests"
	"shop/internal/pkg/custom_error"
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
func (h HomeRepository) GetProduct(ctx context.Context, productSlug, sku string) (entities.Product, error) {
	var product entities.Product
	err := h.db.Where("slug=? and sku=? and status=true", productSlug, sku).First(&product).Error
	return product, err
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
	paginateQuery := pagination.Paginate(c, &products, &pg, h.db)

	if pErr := paginateQuery(h.db).Preload("ProductImages").Where("category_id=?", category.ID).Find(&products).Error; pErr != nil {
		return pg, pErr
	}

	pg.Rows = products
	return pg, nil
}
