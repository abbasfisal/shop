package responses

import (
	"github.com/spf13/viper"
	"shop/internal/entities"
)

type Customer struct {
	ID        uint
	Mobile    string
	FirstName string
	LastName  string
	IsActive  bool
	Address   Address
	Cart      Cart
}
type Address struct {
	ID                 uint
	CustomerID         uint
	ReceiverName       string
	ReceiverMobile     string
	ReceiverAddress    string
	ReceiverPostalCode string
}

func toAddress(address entities.Address) Address {
	return Address{
		ID:                 address.ID,
		CustomerID:         address.CustomerID,
		ReceiverName:       address.ReceiverName,
		ReceiverMobile:     address.ReceiverMobile,
		ReceiverAddress:    address.ReceiverAddress,
		ReceiverPostalCode: address.ReceiverPostalCode,
	}
}

func ToCustomer(cus entities.Customer) Customer {

	return Customer{
		ID:        cus.ID,
		Mobile:    cus.Mobile,
		FirstName: cus.FirstName,
		LastName:  cus.LastName,
		IsActive:  cus.Active,
		Address:   toAddress(cus.Address),
		Cart:      toCart(cus.Carts), // در اینجا اولین سبد خرید مشتری در نظر گرفته می‌شود

	}
}

//-----------

type Cart struct {
	ID         uint
	CustomerID uint
	Status     uint8

	CartItem CartItem
}

type CartItem struct {
	TotalItemCount     uint
	TotalSalePrice     uint
	TotalOriginalPrice uint
	TotalProfitPrice   uint

	Data []Item
}

type Item struct {
	ID            uint
	CustomerID    uint
	CartID        uint
	ProductID     uint
	InventoryID   uint
	Quantity      uint8
	OriginalPrice uint
	SalePrice     uint

	ProductSku   string
	ProductTitle string
	ProductImage string
	ProductSlug  string
}

func toCart(cartEntity []entities.Cart) Cart {
	if len(cartEntity) > 0 {
		return Cart{
			ID:         cartEntity[0].ID,
			CustomerID: cartEntity[0].CustomerID,
			Status:     cartEntity[0].Status,
			CartItem:   toCartItems(cartEntity[0].CartItems), // تبدیل آیتم‌های سبد خرید
		}
	}
	return Cart{}

}

func toCartItems(cartItemEntities []entities.CartItem) CartItem {
	var items []Item
	var totalItemCount uint
	var totalSalePrice, totalOriginalPrice, totalProfitPrice uint

	counter := uint(1)
	// محاسبه مقادیر و تبدیل آیتم‌ها
	for _, item := range cartItemEntities {
		items = append(items, toCartItem(item))
		totalItemCount += counter //----------------
		totalSalePrice += uint(item.SalePrice) * uint(item.Quantity)
		totalOriginalPrice += uint(item.OriginalPrice) * uint(item.Quantity)
		totalProfitPrice += (uint(item.OriginalPrice) - uint(item.SalePrice)) * uint(item.Quantity)

		counter++
	}

	return CartItem{
		TotalItemCount:     totalItemCount,
		TotalSalePrice:     totalSalePrice,
		TotalOriginalPrice: totalOriginalPrice,
		TotalProfitPrice:   totalProfitPrice,
		Data:               items, // تبدیل لیست آیتم‌ها
	}
}

func toCartItem(cartItem entities.CartItem) Item {
	return Item{
		ID:            cartItem.ID,
		CustomerID:    cartItem.CustomerID,
		CartID:        cartItem.CartID,
		ProductID:     cartItem.ProductID,
		InventoryID:   cartItem.InventoryID,
		Quantity:      cartItem.Quantity,
		OriginalPrice: cartItem.OriginalPrice,
		SalePrice:     cartItem.SalePrice,
		ProductSku:    cartItem.ProductSku,
		ProductTitle:  cartItem.ProductTitle,
		ProductImage:  viper.GetString("Upload.Products") + cartItem.ProductImage, // ترکیب مسیر تصویر با پیکربندی
		ProductSlug:   cartItem.ProductSlug,
	}
}
