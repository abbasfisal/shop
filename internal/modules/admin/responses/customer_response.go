package responses

import "shop/internal/entities"

type Customer struct {
	ID        uint
	FirstName string
	LastName  string
	Mobile    string
	Status    bool
}

type Customers struct {
	Data []Customer
}

func ToCustomers(customers []*entities.Customer) *Customers {
	var cResponse Customers
	for _, p := range customers {
		cResponse.Data = append(cResponse.Data, *ToCustomer(p))
	}
	return &cResponse
}
func ToCustomer(c *entities.Customer) *Customer {
	return &Customer{
		ID:        c.ID,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Mobile:    c.Mobile,
		Status:    c.Active,
	}
}
