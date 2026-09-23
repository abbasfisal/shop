package entities

import "time"

type DashboardStats struct {
	TodaySales    float64 `gorm:"column:today_sales"`
	MonthSales    float64 `gorm:"column:month_sales"`
	TotalRevenue  float64 `gorm:"column:total_revenue"`
	TodayOrders   int     `gorm:"column:today_orders"`
	PendingOrders int     `gorm:"column:pending_orders"`
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
