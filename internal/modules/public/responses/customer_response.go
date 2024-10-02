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
	Carts Carts
}
type Cart struct {
	ID            uint
	CustomerID    uint
	ProductID     uint
	InventoryID   uint
	Count         uint8
	Status        uint8
	ProductSku    string
	ProductTitle  string
	ProductSlug   string
	ProductImage  string
	OriginalPrice uint
	SalePrice     uint
}
type Carts struct {
	TotalItemCount     uint
	TotalSalePrice     uint
	TotalOriginalPrice uint
	TotalProfitPrice   uint

	Data []Cart
}

func toCart(cartItem entities.Cart) Cart {
	return Cart{
		ID:          cartItem.ID,
		CustomerID:  cartItem.CustomerID,
		ProductID:   cartItem.ProductID,
		InventoryID: cartItem.InventoryID,

		Count:  cartItem.Count,
		Status: cartItem.Status,

		ProductSku:   cartItem.ProductSku,
		ProductTitle: cartItem.ProductTitle,
		ProductSlug:  cartItem.ProductSlug,
		ProductImage: viper.GetString("Upload.Products") + cartItem.ProductImage,

		OriginalPrice: cartItem.OriginalPrice,
		SalePrice:     cartItem.SalePrice,
	}
}

func toCarts(cartData []entities.Cart) Carts {
	var carts Carts
	counter := 1
	for _, item := range cartData {
		carts.TotalItemCount += uint(counter)
		carts.TotalProfitPrice += (uint(item.OriginalPrice) - uint(item.SalePrice)) * uint(item.Count)

		carts.TotalOriginalPrice += item.OriginalPrice * uint(item.Count)
		carts.TotalSalePrice += item.SalePrice * uint(item.Count)

		carts.Data = append(carts.Data, toCart(item))
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
