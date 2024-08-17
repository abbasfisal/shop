package responses

import (
	"shop/internal/entities"
)

type Customer struct {
	ID        uint
	Mobile    string
	FirstName string
	LastName  string
	IsActive  bool
	//address

}

func ToCustomer(cus entities.Customer) Customer {

	return Customer{
		ID:        cus.ID,
		Mobile:    cus.Mobile,
		FirstName: cus.FirstName,
		LastName:  cus.LastName,
		IsActive:  cus.Active,
	}

}
