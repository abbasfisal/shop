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
	//address
	Carts []Cart
}
type Cart struct {
	ID            uint
	CustomerID    uint
	ProductID     uint
	InventoryID   uint
	Count         uint8
	Status        uint8
	ProductTitle  string
	ProductSlug   string
	ProductImage  string
	OriginalPrice uint
	SalePrice     uint
}

func toCart(cartItem entities.Cart) Cart {
	return Cart{
		ID:            cartItem.ID,
		CustomerID:    cartItem.CustomerID,
		ProductID:     cartItem.ProductID,
		InventoryID:   cartItem.InventoryID,
		Count:         cartItem.Count,
		Status:        cartItem.Status,
		ProductTitle:  cartItem.ProductTitle,
		ProductSlug:   cartItem.ProductSlug,
		ProductImage:  viper.GetString("Upload.Products") + cartItem.ProductImage,
		OriginalPrice: cartItem.OriginalPrice,
		SalePrice:     cartItem.SalePrice,
	}
}

func toCarts(cartData []entities.Cart) []Cart {
	var carts []Cart
	for _, item := range cartData {
		carts = append(carts, toCart(item))
	}
	return carts
}

func ToCustomer(cus entities.Customer) Customer {

	return Customer{
		ID:        cus.ID,
		Mobile:    cus.Mobile,
		FirstName: cus.FirstName,
		LastName:  cus.LastName,
		IsActive:  cus.Active,
		Carts:     toCarts(cus.Carts),
	}

}
