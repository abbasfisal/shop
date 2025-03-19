package dashboard

import (
	"gorm.io/gorm"
	"shop/internal/entities"
	"time"
)

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepositoryInterface {
	return &DashboardRepository{db: db}
}

func (d *DashboardRepository) GetDashboardStates() (*DashboardData, error) {
	var data DashboardData
	today := time.Now().Truncate(24 * time.Hour)
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	// دریافت خلاصه وضعیت
	// some another counts 👇
	//COUNT(CASE WHEN order_status = 'shipped' THEN 1 END) AS unshipped_orders,
	//COUNT(CASE WHEN order_status = 'shipped' THEN 1 END) AS unshipped_orders,
	//COUNT(CASE WHEN order_status = 'shipped' THEN 1 END) AS unshipped_orders

	err := d.db.Raw(`
		SELECT 
			SUM(CASE WHEN created_at >= ? THEN total_sale_price ELSE 0 END) AS today_sales,
			SUM(CASE WHEN created_at >= ? THEN total_sale_price ELSE 0 END) AS month_sales,
			SUM(total_sale_price) AS total_revenue,
			COUNT(CASE WHEN created_at >= ? THEN 1 END) AS today_orders,
			COUNT(CASE WHEN order_status = ? THEN 1 END) AS pending_orders
			
		FROM orders
	`, today, monthStart, today, entities.OrderPending).Scan(&data.Stats).Error
	if err != nil {
		return nil, err
	}

	// دریافت لیست ۱۰ سفارش آخر
	err = d.db.Raw(`
		SELECT o.id, c.first_name AS customer, o.total_sale_price AS total, o.order_status AS status, o.created_at
		FROM orders o
		JOIN customers c ON o.customer_id = c.id
		ORDER BY o.created_at DESC
		LIMIT 10
	`).Scan(&data.RecentOrders).Error
	if err != nil {
		return nil, err
	}

	// دریافت لیست کاربران جدید
	err = d.db.Raw(`
		SELECT c.id, c.first_name, c.mobile, c.created_at, COUNT(o.id) AS order_count
		FROM customers c
		LEFT JOIN orders o ON c.id = o.customer_id
		GROUP BY c.id
		ORDER BY c.created_at DESC
		LIMIT 10
	`).Scan(&data.NewUsers).Error
	if err != nil {
		return nil, err
	}

	// statical customer report
	err = d.db.Raw(`
		SELECT
			COUNT(id) AS total_customer
		FROM customers
	`).Scan(&data.StaticalReport.TotalCustomer).Error
	if err != nil {
		return nil, err
	}

	// statical product report
	err = d.db.Raw(`
		SELECT
			COUNT(id) AS total_product
		FROM products
	`).Scan(&data.StaticalReport.TotalProduct).Error
	if err != nil {
		return nil, err
	}

	//// دریافت کالاهای با موجودی کم
	//err = d.db.Raw(`
	//	SELECT id, name, stock
	//	FROM products
	//	WHERE stock < 10
	//	ORDER BY stock ASC
	//	LIMIT 10
	//`).Scan(&data.LowStockItems).Error
	//if err != nil {
	//	return nil, err
	//}

	// دریافت گزارش پرداخت‌ها
	err = d.db.Raw(`
		SELECT 
			COUNT(CASE WHEN status = ? THEN 1 END) AS successful_payments,
			COUNT(CASE WHEN status = ? THEN 1 END) AS failed_payments,
			SUM(CASE WHEN status = ? THEN amount ELSE 0 END) AS total_successful_sales
		FROM payments
	`, entities.OrderConfirmed, entities.OrderCancelled, entities.OrderConfirmed).Scan(&data.PaymentReport).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

type DashboardStats struct {
	TodaySales      float64 `gorm:"column:today_sales"`
	MonthSales      float64 `gorm:"column:month_sales"`
	TotalRevenue    float64 `gorm:"column:total_revenue"`
	TodayOrders     int     `gorm:"column:today_orders"`
	PendingOrders   int     `gorm:"column:pending_orders"`
	UnshippedOrders int     `gorm:"column:unshipped_orders"`
}

type RecentOrder struct {
	ID        uint      `gorm:"column:id"`
	Customer  string    `gorm:"column:customer"`
	Total     float64   `gorm:"column:total"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type NewUser struct {
	ID         uint      `gorm:"column:id"`
	FirstName  string    `gorm:"column:first_name"`
	Mobile     string    `gorm:"column:mobile"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	OrderCount int       `gorm:"column:order_count"`
}

type LowStockProduct struct {
	ID    uint   `gorm:"column:id"`
	Name  string `gorm:"column:name"`
	Stock int    `gorm:"column:stock"`
}

type PaymentReport struct {
	SuccessfulPayments   int     `gorm:"column:successful_payments"`
	FailedPayments       int     `gorm:"column:failed_payments"`
	TotalSuccessfulSales float64 `gorm:"column:total_successful_sales"`
}
type StaticalReport struct {
	TotalCustomer int `gorm:"column:total_customer"`
	TotalProduct  int `gorm:"column:total_product"`
}

type DashboardData struct {
	Stats          DashboardStats
	RecentOrders   []RecentOrder
	NewUsers       []NewUser
	LowStockItems  []LowStockProduct
	PaymentReport  PaymentReport
	StaticalReport StaticalReport
}
